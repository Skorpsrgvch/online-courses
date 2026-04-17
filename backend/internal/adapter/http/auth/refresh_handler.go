package auth

import (
	"net/http"
	"strings"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/gin-gonic/gin"
)

type RefreshHandler struct{}

func NewRefreshHandler() *RefreshHandler {
	return &RefreshHandler{}
}

func (h *RefreshHandler) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		common.HandleError(c, common.HttpError("missing refresh token", http.StatusUnauthorized))
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := middleware.ParseToken(tokenString)
	if err != nil {
		common.HandleError(c, common.HttpError("invalid refresh token", http.StatusUnauthorized))
		return
	}

	// Генерируем новую пару токенов
	accessToken, refreshToken, err := middleware.GenerateToken(
		claims.UserID,
		claims.Email,
		claims.Name,
		claims.Role,
	)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    900,
	})
}
