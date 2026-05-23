package service

import (
	"net/http"

	listUC "github.com/Skorpsrgvch/online-courses/internal/usecase/service/list"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ListHandler struct {
	usecase *listUC.Usecase
}

func NewListHandler(usecase *listUC.Usecase) *ListHandler {
	return &ListHandler{usecase: usecase}
}

func (h *ListHandler) Handle(c *gin.Context) {
	zap.L().Debug("Fetching all services")

	services, err := h.usecase.Execute(c.Request.Context())
	if err != nil {
		zap.L().Error("Failed to list services", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	zap.L().Info("Services listed successfully", zap.Int("count", len(services)))
	c.JSON(http.StatusOK, services)
}
