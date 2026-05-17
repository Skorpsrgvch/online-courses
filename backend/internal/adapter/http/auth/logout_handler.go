package auth

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/logout"
	"github.com/gin-gonic/gin"
)

type LogoutHandler struct {
	usecase *logout.Usecase
}

func NewLogoutHandler(usecase *logout.Usecase) *LogoutHandler {
	return &LogoutHandler{usecase: usecase}
}

func (h *LogoutHandler) Handle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, common.HttpError("unauthorized", http.StatusUnauthorized))
		return
	}

	if err := h.usecase.Execute(c.Request.Context(), userID); err != nil {
		common.HandleError(c, err)
		return
	}

	// Удаляем куку с refresh-токеном на клиенте
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
