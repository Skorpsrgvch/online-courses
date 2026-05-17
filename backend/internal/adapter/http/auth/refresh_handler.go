package auth

import (
	"net/http"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/refresh"
	"github.com/gin-gonic/gin"
)

type RefreshHandler struct {
	usecase *refresh.Usecase
}

// Конструктор теперь принимает usecase
func NewRefreshHandler(usecase *refresh.Usecase) *RefreshHandler {
	return &RefreshHandler{usecase: usecase}
}

func (h *RefreshHandler) Handle(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		common.HandleError(c, common.HttpError("refresh token not found in cookie", http.StatusUnauthorized))
		return
	}

	output, err := h.usecase.Execute(c.Request.Context(), cookie)
	if err != nil {
		common.HandleError(c, common.HttpError(err.Error(), http.StatusUnauthorized))
		return
	}

	c.SetCookie(
		"refresh_token",
		output.RefreshToken,
		int(7*24*time.Hour.Seconds()), // Срок жизни куки
		"/",
		"",
		false, // ПОМЕНЯТЬ ПОТОМ
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": output.AccessToken,
		"expires_in":   900,
	})
}
