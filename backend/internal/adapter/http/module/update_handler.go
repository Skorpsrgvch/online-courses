package module

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	updateUC "github.com/Skorpsrgvch/online-courses/internal/usecase/module/update"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type updateModuleRequest struct {
	Title      string `json:"title" binding:"required"`
	Order      int    `json:"order"`
	WeekNumber int    `json:"week_number"` // Добавлено
}

type UpdateHandler struct {
	usecase *updateUC.Usecase
}

func NewUpdateHandler(usecase *updateUC.Usecase) *UpdateHandler {
	return &UpdateHandler{usecase: usecase}
}

func (h *UpdateHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	moduleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zap.L().Debug("Invalid module ID format", zap.String("id", c.Param("id")), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid module ID", http.StatusBadRequest))
		return
	}

	var req updateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Debug("Invalid JSON in update module", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Debug("Updating module", zap.Int("moduleID", moduleID), zap.String("title", req.Title))

	input := updateUC.Input{
		ID:         moduleID,
		Title:      req.Title,
		Order:      req.Order,
		WeekNumber: req.WeekNumber,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Failed to update module", zap.Int("moduleID", moduleID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Module updated successfully", zap.Int("moduleID", moduleID))
	c.Status(http.StatusOK)
}
