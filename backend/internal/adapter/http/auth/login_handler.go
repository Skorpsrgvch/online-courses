package auth

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	loginUC "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/login"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginHandler struct {
	usecase      *loginUC.Usecase
	refreshStore refreshTokenStore
}

type refreshTokenStore interface {
	Save(ctx context.Context, userID int, token string, expiresAt time.Time) error
}

func NewLoginHandler(usecase *loginUC.Usecase, refreshStore refreshTokenStore) *LoginHandler {
	return &LoginHandler{usecase: usecase, refreshStore: refreshStore}
}

func (h *LoginHandler) Handle(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	zap.L().Debug("Login attempt", zap.String("email", req.Email))

	input := loginUC.Input{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		// Не логируем детали ошибки (неверный пароль) для безопасности,
		// но сам факт неудачи можно залогировать на уровне INFO/WARN без пароля
		zap.L().Info("Login failed", zap.String("email", req.Email), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	accessToken, refreshToken, err := middleware.GenerateToken(
		output.User.ID,
		output.User.Email,
		output.User.Name,
		output.User.Role,
	)
	if err != nil {
		zap.L().Error("Token generation failed", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	if h.refreshStore != nil {
		if err := h.refreshStore.Save(c.Request.Context(), output.User.ID, refreshToken, time.Now().UTC().Add(7*24*time.Hour)); err != nil {
			zap.L().Error("Failed to persist refresh token", zap.Error(err))
			common.HandleError(c, err)
			return
		}
	}

	secureCookie := strings.HasPrefix(strings.ToLower(os.Getenv("FRONTEND_URL")), "https://")
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteLaxMode,
	})

	zap.L().Info("Login successful", zap.Int("userID", output.User.ID), zap.String("email", req.Email))

	c.JSON(http.StatusOK, gin.H{
		"user_id":       output.User.ID,
		"email":         output.User.Email,
		"name":          output.User.Name,
		"role":          output.User.Role,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    900,
	})
}
