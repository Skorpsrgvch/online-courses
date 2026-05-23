package service

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	updateUC "github.com/Skorpsrgvch/online-courses/internal/usecase/service/update"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UpdateHandler struct {
	usecase *updateUC.Usecase
}

func NewUpdateHandler(usecase *updateUC.Usecase) *UpdateHandler {
	return &UpdateHandler{usecase: usecase}
}

func (h *UpdateHandler) Handle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		zap.L().Debug("Invalid service ID format", zap.String("id", idStr), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid service ID", http.StatusBadRequest))
		return
	}

	var input updateUC.Input
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Debug("Invalid JSON in update service request", zap.Error(err))
		common.HandleError(c, common.HttpError("invalid request body", http.StatusBadRequest))
		return
	}
	input.ID = id

	zap.L().Info("Updating service", zap.Int("id", id), zap.String("title", input.Title))

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Failed to update service", zap.Int("id", id), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Service updated successfully", zap.Int("id", id))
	c.Status(http.StatusOK)
}
