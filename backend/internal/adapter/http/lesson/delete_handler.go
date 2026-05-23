package lesson

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	deleteUC "github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/delete"
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
		common.HandleError(c, common.HttpError("admin access required", http.StatusForbidden))
		return
	}

	idStr := c.Param("id")
	lessonID, err := strconv.Atoi(idStr)
	if err != nil {
		zap.L().Debug("Invalid lesson ID format", zap.String("id", idStr), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid lesson ID", http.StatusBadRequest))
		return
	}

	zap.L().Info("Deleting lesson", zap.Int("lessonID", lessonID))

	input := deleteUC.Input{ID: lessonID}
	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Failed to delete lesson", zap.Int("lessonID", lessonID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Lesson deleted successfully", zap.Int("lessonID", lessonID))
	c.Status(http.StatusNoContent)
}
