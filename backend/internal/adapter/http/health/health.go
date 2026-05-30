package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	dbPing func() error // интерфейс для проверки БД
}

func NewHandler(dbPing func() error) *Handler {
	return &Handler{dbPing: dbPing}
}

func (h *Handler) Handle(c *gin.Context) {
	if h.dbPing != nil {
		if err := h.dbPing(); err != nil {
			zap.L().Warn("Health check: DB not reachable", zap.Error(err))
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "database connection failed",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	})
}
