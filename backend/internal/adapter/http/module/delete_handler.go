package module

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	deleteUC "github.com/Skorpsrgvch/online-courses/internal/usecase/module/delete"
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

	zap.L().Debug("Deleting module", zap.Int("moduleID", moduleID))

	input := deleteUC.Input{ID: moduleID}
	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Failed to delete module", zap.Int("moduleID", moduleID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Module deleted successfully", zap.Int("moduleID", moduleID))
	c.Status(http.StatusNoContent)
}
