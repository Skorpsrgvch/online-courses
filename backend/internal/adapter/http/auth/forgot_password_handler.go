package auth

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/forgotpassword"
	"github.com/gin-gonic/gin"
)

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordHandler struct {
	usecase *forgotpassword.Usecase
}

func NewForgotPasswordHandler(usecase *forgotpassword.Usecase) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{usecase: usecase}
}

func (h *ForgotPasswordHandler) Handle(c *gin.Context) {
	var req forgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	input := forgotpassword.Input{Email: req.Email}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
