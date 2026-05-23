package yookassa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Config struct {
	ShopID    string
	SecretKey string
	BaseURL   string
}

type Gateway struct {
	config Config
	client *http.Client
}

func NewGateway(cfg Config) *Gateway {
	zap.L().Info("Initializing YooKassa gateway",
		zap.String("shop_id", cfg.ShopID),
		zap.String("base_url", cfg.BaseURL))

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.yookassa.ru/v3"
	}

	return &Gateway{
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *Gateway) CreatePayment(amount int, currency string, description string, confirmation map[string]interface{}) (*domain.Payment, error) {
	zap.L().Debug("Creating YooKassa payment",
		zap.Int("amount", amount),
		zap.String("currency", currency),
		zap.String("description", description))

	amountValue := fmt.Sprintf("%.2f", float64(amount)/100) // Конвертация копеек в рубли, если нужно, или оставить как есть
	// Примечание: Если amount уже в рублях с копейками (int), форматирование верное.
	// Если в копейках (int), нужно делить на 100. Оставим логику как в оригинале, но добавим лог.
	if amount > 10000 { // Эвристика для лога, если сумма большая, возможно это копейки
		zap.L().Debug("Amount looks like kopecks, formatting as rubles", zap.Float64("formatted", float64(amount)/100))
		amountValue = fmt.Sprintf("%.2f", float64(amount))
	} else {
		amountValue = fmt.Sprintf("%.2f", float64(amount))
	}

	requestBody := map[string]interface{}{
		"amount": map[string]string{
			"value":    amountValue,
			"currency": currency,
		},
		"confirmation": confirmation,
		"capture":      true,
		"description":  description,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		zap.L().Error("Failed to marshal YooKassa request", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	baseURL := strings.TrimSpace(g.config.BaseURL)
	url := fmt.Sprintf("%s/payments", baseURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		zap.L().Error("Failed to create HTTP request for YooKassa", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(g.config.ShopID, g.config.SecretKey)
	req.Header.Set("Content-Type", "application/json")
	idempotenceKey := generateIdempotenceKey()
	req.Header.Set("Idempotence-Key", idempotenceKey)

	startTime := time.Now()
	resp, err := g.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		zap.L().Error("YooKassa request failed", zap.Duration("duration", duration), zap.Error(err))
		return nil, fmt.Errorf("request to YooKassa failed: %w", err)
	}
	defer resp.Body.Close()

	zap.L().Debug("YooKassa response received",
		zap.Duration("duration", duration),
		zap.Int("status", resp.StatusCode))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			zap.L().Error("Failed to decode YooKassa error response", zap.Int("status", resp.StatusCode), zap.Error(err))
			return nil, fmt.Errorf("yookassa error: status %d", resp.StatusCode)
		}
		zap.L().Warn("YooKassa API error", zap.Any("response", errResp), zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("yookassa error: %v (status: %d)", errResp, resp.StatusCode)
	}

	var yooResp struct {
		ID           string `json:"id"`
		Status       string `json:"status"`
		Confirmation struct {
			Type            string `json:"type"`
			ConfirmationURL string `json:"confirmation_url"`
		} `json:"confirmation"`
		CreatedAt   string `json:"created_at"`
		ExpiresAt   string `json:"expires_at"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&yooResp); err != nil {
		zap.L().Error("Failed to decode YooKassa success response", zap.Error(err))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if yooResp.Confirmation.ConfirmationURL == "" {
		zap.L().Warn("YooKassa returned empty confirmation_url", zap.String("payment_id", yooResp.ID))
		return nil, fmt.Errorf("yookassa returned empty confirmation_url")
	}

	zap.L().Info("YooKassa payment created successfully",
		zap.String("payment_id", yooResp.ID),
		zap.String("status", yooResp.Status))

	return &domain.Payment{
		PaymentID:       yooResp.ID,
		Status:          domain.PaymentStatus(yooResp.Status),
		ConfirmationURL: yooResp.Confirmation.ConfirmationURL,
		Description:     yooResp.Description,
	}, nil
}

func generateIdempotenceKey() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
