package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	logoutUC "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/logout"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LogoutHandler struct {
	usecase *logoutUC.Usecase
}

func NewLogoutHandler(usecase *logoutUC.Usecase) *LogoutHandler {
	return &LogoutHandler{usecase: usecase}
}

func (h *LogoutHandler) Handle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, common.HttpError("unauthorized", http.StatusUnauthorized))
		return
	}

	zap.L().Debug("Logout request", zap.Int("userID", userID))

	if err := h.usecase.Execute(c.Request.Context(), userID); err != nil {
		zap.L().Error("Logout execution failed", zap.Int("userID", userID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	// Удаляем куку
	secureCookie := strings.HasPrefix(strings.ToLower(os.Getenv("FRONTEND_URL")), "https://")
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteLaxMode,
	})

	zap.L().Info("User logged out", zap.Int("userID", userID))
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
