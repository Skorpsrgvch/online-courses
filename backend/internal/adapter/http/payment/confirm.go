package payment

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/payment/confirm"
	"github.com/gin-gonic/gin"
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

	input := confirm.Input{
		PaymentID: paymentID,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// HandleGet обрабатывает GET-запрос (возврат пользователя после оплаты)
func (h *ConfirmHandler) HandleGet(c *gin.Context) {
	paymentID := c.Param("payment_id")
	returnURL := c.Query("return_url")

	if paymentID == "" {
		common.HandleError(c, common.HttpError("payment_id is required", http.StatusBadRequest))
		return
	}

	input := confirm.Input{
		PaymentID: paymentID,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {

		redirectURL := "/courses?error=payment_failed"
		if returnURL != "" {

			if len(returnURL) > 0 && returnURL[0] == '/' {
				redirectURL = returnURL + "?error=payment_failed"
			}
		}
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Если успешно, редиректим на страницу успеха
	if output.Success && returnURL != "" {
		c.Redirect(http.StatusFound, returnURL)
		return
	}

	c.Redirect(http.StatusFound, "/courses?error=unknown_status")
}
