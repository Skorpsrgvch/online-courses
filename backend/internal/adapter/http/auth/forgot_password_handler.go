package auth

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	forgotPassUC "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/forgotpassword"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordHandler struct {
	usecase *forgotPassUC.Usecase
}

func NewForgotPasswordHandler(usecase *forgotPassUC.Usecase) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{usecase: usecase}
}

func (h *ForgotPasswordHandler) Handle(c *gin.Context) {
	var req forgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	zap.L().Debug("Forgot password request", zap.String("email", req.Email))

	input := forgotPassUC.Input{Email: req.Email}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Forgot password execution failed", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
