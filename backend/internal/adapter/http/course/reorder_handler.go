package course

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	updateFullUC "github.com/Skorpsrgvch/online-courses/internal/usecase/course/updatefullcourse"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ReorderHandler struct {
	usecase *updateFullUC.Usecase
}

func NewReorderHandler(usecase *updateFullUC.Usecase) *ReorderHandler {
	return &ReorderHandler{usecase: usecase}
}

type ReorderModulesRequest struct {
	ModuleIDs []int `json:"module_ids" binding:"required,min=1"`
}

func (h *ReorderHandler) HandleModules(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, domain.ErrInvalidID)
		return
	}

	var req ReorderModulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Reorder modules request", zap.Int("courseID", courseID), zap.Int("count", len(req.ModuleIDs)))

	if err := h.usecase.ReorderModules(c.Request.Context(), courseID, req.ModuleIDs); err != nil {
		zap.L().Error("Reorder modules failed", zap.Int("courseID", courseID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Modules reordered successfully", zap.Int("courseID", courseID))
	c.Status(http.StatusOK)
}

type ReorderLessonsRequest struct {
	LessonIDs []int `json:"lesson_ids" binding:"required,min=1"`
}

func (h *ReorderHandler) HandleLessons(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	moduleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, domain.ErrInvalidID)
		return
	}

	var req ReorderLessonsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Reorder lessons request", zap.Int("moduleID", moduleID), zap.Int("count", len(req.LessonIDs)))

	if err := h.usecase.ReorderLessons(c.Request.Context(), moduleID, req.LessonIDs); err != nil {
		zap.L().Error("Reorder lessons failed", zap.Int("moduleID", moduleID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Lessons reordered successfully", zap.Int("moduleID", moduleID))
	c.Status(http.StatusOK)
}
