package course

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	deleteUC "github.com/Skorpsrgvch/online-courses/internal/usecase/course/delete"
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

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, common.HttpError("invalid course ID", http.StatusBadRequest))
		return
	}

	zap.L().Info("Delete course request", zap.Int("courseID", courseID))

	input := deleteUC.Input{CourseID: courseID}
	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Delete course failed", zap.Int("courseID", courseID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Course deleted successfully", zap.Int("courseID", courseID))
	c.Status(http.StatusNoContent)
}
