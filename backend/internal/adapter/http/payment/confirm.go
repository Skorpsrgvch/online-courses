package payment

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/payment/confirm"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ConfirmHandler struct {
	usecase *confirm.UseCase
}

func NewConfirmHandler(usecase *confirm.UseCase) *ConfirmHandler {
	return &ConfirmHandler{usecase: usecase}
}

func (h *ConfirmHandler) Handle(c *gin.Context) {
	paymentID := c.Param("payment_id")
	if paymentID == "" {
		common.HandleError(c, common.HttpError("payment_id is required", http.StatusBadRequest))
		return
	}

	zap.L().Debug("Confirm payment request (POST)", zap.String("payment_id", paymentID))

	input := confirm.Input{PaymentID: paymentID}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *ConfirmHandler) HandleGet(c *gin.Context) {
	paymentID := c.Param("payment_id")
	returnURL := c.Query("return_url")

	if paymentID == "" {
		common.HandleError(c, common.HttpError("payment_id is required", http.StatusBadRequest))
		return
	}

	zap.L().Debug("Confirm payment request (GET/Return)",
		zap.String("payment_id", paymentID),
		zap.String("return_url", returnURL))

	input := confirm.Input{PaymentID: paymentID}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Warn("Payment confirmation failed on return", zap.Error(err))
		redirectURL := "/courses?error=payment_failed"
		if returnURL != "" && len(returnURL) > 0 && returnURL[0] == '/' {
			redirectURL = returnURL + "?error=payment_failed"
		}
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	if output.Success && returnURL != "" {
		zap.L().Info("Payment confirmed successfully, redirecting", zap.String("url", returnURL))
		c.Redirect(http.StatusFound, returnURL)
		return
	}

	zap.L().Warn("Payment status unknown or failed, redirecting to default", zap.Bool("success", output.Success))
	c.Redirect(http.StatusFound, "/courses?error=unknown_status")
}
