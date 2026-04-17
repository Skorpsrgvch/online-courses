package auth

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
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

	// Генерируем JWT
	accessToken, refreshToken, err := middleware.GenerateToken(
		output.User.ID,
		output.User.Email,
		output.User.Name,
		output.User.Role,
	)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       output.User.ID,
		"email":         output.User.Email,
		"name":          output.User.Name,
		"role":          output.User.Role,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    900, // 15 минут в секундах
	})
}
