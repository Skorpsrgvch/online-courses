// internal/adapter/http/auth/login_handler.go
package auth

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/login"
	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginHandler struct {
	usecase *login.Usecase
}

func NewLoginHandler(usecase *login.Usecase) *LoginHandler {
	return &LoginHandler{usecase: usecase}
}

func (h *LoginHandler) Handle(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	input := login.Input{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Для диплома можно вернуть ID + роль
	c.JSON(http.StatusOK, gin.H{
		"user_id": output.User.ID,
		"role":    output.User.Role,
	})
}
