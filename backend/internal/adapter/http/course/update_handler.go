package course

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	updateUC "github.com/Skorpsrgvch/online-courses/internal/usecase/course/update"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type updateCourseRequest struct {
	Title             string             `json:"title" binding:"required"`
	Description       string             `json:"description"`
	IsPublic          bool               `json:"is_public"`
	Price             int                `json:"price"`
	IsActive          bool               `json:"is_active"`
	CoverImageURL     string             `json:"cover_image_url"`
	Contraindications string             `json:"contraindications"`
	Recommendations   string             `json:"recommendations"`
	TargetAudience    string             `json:"target_audience"`
	CourseBasis       string             `json:"course_basis"`
	ClassBasis        string             `json:"class_basis"`
	Bonuses           []domain.BonusItem `json:"bonuses"`
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

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zap.L().Debug("Invalid course ID format", zap.Error(err))
		common.HandleError(c, domain.ErrInvalidID)
		return
	}

	var req updateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Update course request", zap.Int("courseID", courseID))

	course := &domain.Course{
		ID:                courseID,
		Title:             req.Title,
		Description:       req.Description,
		IsPublic:          req.IsPublic,
		Price:             req.Price,
		IsActive:          req.IsActive,
		CoverImageURL:     req.CoverImageURL,
		Contraindications: req.Contraindications,
		Recommendations:   req.Recommendations,
		TargetAudience:    req.TargetAudience,
		CourseBasis:       req.CourseBasis,
		ClassBasis:        req.ClassBasis,
		Bonuses:           req.Bonuses,
	}

	if err := h.usecase.Execute(c.Request.Context(), course); err != nil {
		zap.L().Error("Update course failed", zap.Int("courseID", courseID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Course updated successfully", zap.Int("courseID", courseID))
	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

func (h *UpdateHandler) HandleStatusPatch(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, domain.ErrInvalidID)
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Update course status request", zap.Int("courseID", courseID), zap.Bool("isActive", req.IsActive))

	if err := h.usecase.UpdateStatus(c.Request.Context(), courseID, req.IsActive); err != nil {
		zap.L().Error("Update course status failed", zap.Int("courseID", courseID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Course status updated successfully", zap.Int("courseID", courseID))
	c.Status(http.StatusOK)
}
