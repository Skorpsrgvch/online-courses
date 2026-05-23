package auth

import (
	"net/http"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	refreshUC "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/refresh"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RefreshHandler struct {
	usecase *refreshUC.Usecase
}

func NewRefreshHandler(usecase *refreshUC.Usecase) *RefreshHandler {
	return &RefreshHandler{usecase: usecase}
}

func (h *RefreshHandler) Handle(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		zap.L().Debug("Refresh token missing in cookie")
		common.HandleError(c, common.HttpError("refresh token not found", http.StatusUnauthorized))
		return
	}

	zap.L().Debug("Refresh token request received")

	output, err := h.usecase.Execute(c.Request.Context(), cookie)
	if err != nil {
		zap.L().Info("Refresh token validation failed", zap.Error(err))
		common.HandleError(c, common.HttpError(err.Error(), http.StatusUnauthorized))
		return
	}

	c.SetCookie(
		"refresh_token",
		output.RefreshToken,
		int(7*24*time.Hour.Seconds()),
		"/",
		"",
		false, // Secure: false для локальной разработки, true для prod (HTTPS)
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": output.AccessToken,
		"expires_in":   900,
	})
}
