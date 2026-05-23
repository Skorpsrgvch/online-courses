package payment

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/usecase/payment/callback"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WebhookRequest struct {
	Type   string `json:"type"`
	Event  string `json:"event"`
	Object struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"object"`
}

type CallbackHandler struct {
	usecase *callback.UseCase
}

func NewCallbackHandler(usecase *callback.UseCase) *CallbackHandler {
	return &CallbackHandler{usecase: usecase}
}

func (h *CallbackHandler) Handle(c *gin.Context) {
	var req WebhookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Warn("Invalid JSON in payment callback", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	paymentID := req.Object.ID
	status := req.Object.Status

	zap.L().Info("Payment callback received",
		zap.String("payment_id", paymentID),
		zap.String("status", status),
		zap.String("event_type", req.Type))

	if paymentID == "" {
		zap.L().Warn("Callback missing payment_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment id is missing"})
		return
	}

	input := callback.Input{
		PaymentID: paymentID,
		Status:    status,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Payment callback processing failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
