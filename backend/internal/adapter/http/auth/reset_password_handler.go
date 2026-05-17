package auth

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/resetpassword"
	"github.com/gin-gonic/gin"
)

type resetPasswordRequest struct {
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ResetPasswordHandler struct {
	usecase *resetpassword.Usecase
}

func NewResetPasswordHandler(usecase *resetpassword.Usecase) *ResetPasswordHandler {
	return &ResetPasswordHandler{usecase: usecase}
}

func (h *ResetPasswordHandler) Handle(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	input := resetpassword.Input{
		Code:        req.Code,
		NewPassword: req.NewPassword,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}
