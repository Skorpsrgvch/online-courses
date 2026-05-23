package auth

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	resetPassUC "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/resetpassword"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type resetPasswordRequest struct {
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ResetPasswordHandler struct {
	usecase *resetPassUC.Usecase
}

func NewResetPasswordHandler(usecase *resetPassUC.Usecase) *ResetPasswordHandler {
	return &ResetPasswordHandler{usecase: usecase}
}

func (h *ResetPasswordHandler) Handle(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	zap.L().Debug("Reset password request received")

	input := resetPassUC.Input{
		Code:        req.Code,
		NewPassword: req.NewPassword,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Info("Reset password failed", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Password reset successful")
	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}
