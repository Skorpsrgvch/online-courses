package course

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	updatefullcourse "github.com/Skorpsrgvch/online-courses/internal/usecase/course/updatefullcourse"
	"github.com/gin-gonic/gin"
)

type ReorderHandler struct {
	usecase *updatefullcourse.Usecase
}

func NewReorderHandler(usecase *updatefullcourse.Usecase) *ReorderHandler {
	return &ReorderHandler{usecase: usecase}
}

// ReorderModulesRequest ожидает массив ID в новом порядке
type ReorderModulesRequest struct {
	ModuleIDs []int `json:"module_ids" binding:"required,min=1"`
}

func (h *ReorderHandler) HandleModules(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		common.HandleError(c, domain.ErrInvalidID)
		return
	}

	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	var req ReorderModulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	if err := h.usecase.ReorderModules(c.Request.Context(), courseID, req.ModuleIDs); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// ReorderLessonsRequest
type ReorderLessonsRequest struct {
	LessonIDs []int `json:"lesson_ids" binding:"required,min=1"`
}

func (h *ReorderHandler) HandleLessons(c *gin.Context) {
	moduleIDStr := c.Param("id")
	moduleID, err := strconv.Atoi(moduleIDStr)
	if err != nil {
		common.HandleError(c, domain.ErrInvalidID)
		return
	}

	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	var req ReorderLessonsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	if err := h.usecase.ReorderLessons(c.Request.Context(), moduleID, req.LessonIDs); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
