package yookassa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

// Config конфигурация шлюза
type Config struct {
	ShopID    string
	SecretKey string
	BaseURL   string
}

// Gateway реализует интерфейс PaymentGateway
type Gateway struct {
	config Config
	client *http.Client
}

func NewGateway(cfg Config) *Gateway {
	// Логирование инициализации (скрытие секретного ключа)
	log.Printf("[YooKassa] Initializing gateway with ShopID: %s, BaseURL: %s", cfg.ShopID, cfg.BaseURL)

	return &Gateway{
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// CreatePayment создает платеж в ЮKassa
func (g *Gateway) CreatePayment(amount int, currency string, description string, confirmation map[string]interface{}) (*domain.Payment, error) {
	log.Printf("[YooKassa] Creating payment: Amount=%d (%s), Desc=%s", amount, currency, description)

	// 1. Форматирование суммы
	amountValue := fmt.Sprintf("%.2f", float64(amount))
	log.Printf("[YooKassa] Formatted amount value: %s", amountValue)

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
		log.Printf("[YooKassa] ERROR marshaling request: %v", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	log.Printf("[YooKassa] Request JSON: %s", string(jsonData))

	// 2. Подготовка URL
	baseURL := strings.TrimSpace(g.config.BaseURL)
	if baseURL == "" {
		baseURL = "https://api.yookassa.ru/v3"
		log.Printf("[YooKassa] Using default BaseURL: %s", baseURL)
	}
	url := fmt.Sprintf("%s/payments", baseURL)
	log.Printf("[YooKassa] Target URL: %s", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[YooKassa] ERROR creating request object: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 3. Настройка заголовков и авторизации
	req.SetBasicAuth(g.config.ShopID, g.config.SecretKey)
	req.Header.Set("Content-Type", "application/json")
	idempotenceKey := generateIdempotenceKey()
	req.Header.Set("Idempotence-Key", idempotenceKey)

	log.Printf("[YooKassa] Sending request with Idempotence-Key: %s", idempotenceKey)

	// 4. Отправка запроса
	startTime := time.Now()
	resp, err := g.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		log.Printf("[YooKassa] ERROR performing HTTP request after %v: %v", duration, err)
		return nil, fmt.Errorf("request to YooKassa failed: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("[YooKassa] Response received in %v with Status: %s", duration, resp.Status)

	// 5. Обработка статус кода
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Printf("[YooKassa] ERROR decoding error response body: %v", err)
			return nil, fmt.Errorf("yookassa error: status %d, failed to decode error response", resp.StatusCode)
		}
		log.Printf("[YooKassa] API Error Response: %v", errResp)
		return nil, fmt.Errorf("yookassa error: %v (status: %d)", errResp, resp.StatusCode)
	}

	// 6. Парсинг успешного ответа
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
		log.Printf("[YooKassa] ERROR decoding success response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("[YooKassa] Payment created successfully. ID: %s, Status: %s, URL: %s",
		yooResp.ID, yooResp.Status, yooResp.Confirmation.ConfirmationURL)

	if yooResp.Confirmation.ConfirmationURL == "" {
		log.Printf("[YooKassa] WARNING: Empty confirmation_url in response")
		return nil, fmt.Errorf("yookassa returned empty confirmation_url")
	}

	payment := &domain.Payment{
		PaymentID:       yooResp.ID,
		Status:          domain.PaymentStatus(yooResp.Status),
		ConfirmationURL: yooResp.Confirmation.ConfirmationURL,
		Description:     yooResp.Description,
	}

	return payment, nil
}

func generateIdempotenceKey() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
