package course

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	listUC "github.com/Skorpsrgvch/online-courses/internal/usecase/course/list"
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
	zap.L().Debug("List courses request (public)")

	output, err := h.usecase.Execute(c.Request.Context())
	if err != nil {
		zap.L().Error("List courses failed", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("List courses successful", zap.Int("count", len(output.Courses)))
	c.JSON(http.StatusOK, output.Courses)
}

func (h *ListHandler) HandleAdmin(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	zap.L().Debug("List courses request (admin)")

	courses, err := h.usecase.ExecuteAdmin(c.Request.Context())
	if err != nil {
		zap.L().Error("List courses (admin) failed", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("List courses (admin) successful", zap.Int("count", len(courses)))
	c.JSON(http.StatusOK, gin.H{"courses": courses})
}
