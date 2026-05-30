package module

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	createUC "github.com/Skorpsrgvch/online-courses/internal/usecase/module/create"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createModuleRequest struct {
	CourseID   int    `json:"course_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Order      int    `json:"order"`
	WeekNumber int    `json:"week_number"`
}

type CreateHandler struct {
	usecase *createUC.Usecase
}

func NewCreateHandler(usecase *createUC.Usecase) *CreateHandler {
	return &CreateHandler{usecase: usecase}
}

func (h *CreateHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	var req createModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Debug("Invalid JSON in create module", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Debug("Creating module", zap.Int("courseID", req.CourseID), zap.String("title", req.Title))

	input := createUC.Input{
		CourseID:   req.CourseID,
		Title:      req.Title,
		Order:      req.Order,
		WeekNumber: req.WeekNumber,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Failed to create module", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Module created successfully", zap.Int("courseID", req.CourseID), zap.String("title", req.Title))
	c.Status(http.StatusCreated)
}
