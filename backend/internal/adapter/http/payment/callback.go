package payment

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/usecase/payment/callback"
	"github.com/gin-gonic/gin"
)

// WebhookRequest структура входящего запроса от ЮKassa
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	paymentID := req.Object.ID
	status := req.Object.Status

	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment id is missing"})
		return
	}

	input := callback.Input{
		PaymentID: paymentID,
		Status:    status,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
