package auth

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/gin-gonic/gin"
)

type MeHandler struct{}

func NewMeHandler() *MeHandler {
	return &MeHandler{}
}

func (h *MeHandler) Handle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, domain.ErrUnauthorized)
		return
	}

	role, _ := c.Get(middleware.RoleKey)
	email, _ := c.Get(middleware.EmailKey)
	name, _ := c.Get(middleware.NameKey)

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   email,
		"name":    name,
		"role":    role,
	})
}
