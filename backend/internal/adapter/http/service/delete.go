package service

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	deleteUC "github.com/Skorpsrgvch/online-courses/internal/usecase/service/delete"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeleteHandler struct {
	usecase *deleteUC.Usecase
}

func NewDeleteHandler(usecase *deleteUC.Usecase) *DeleteHandler {
	return &DeleteHandler{usecase: usecase}
}

func (h *DeleteHandler) Handle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		zap.L().Debug("Invalid service ID format", zap.String("id", idStr), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid service ID", http.StatusBadRequest))
		return
	}

	zap.L().Debug("Deleting service", zap.Int("id", id))

	input := deleteUC.Input{ID: id}
	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Failed to delete service", zap.Int("id", id), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Service deleted successfully", zap.Int("id", id))
	c.Status(http.StatusNoContent)
}
