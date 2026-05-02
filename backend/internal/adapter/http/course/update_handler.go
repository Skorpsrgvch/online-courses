package course

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/update"
	"github.com/gin-gonic/gin"
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
	usecase *update.Usecase
}

func NewUpdateHandler(usecase *update.Usecase) *UpdateHandler {
	return &UpdateHandler{usecase: usecase}
}

// Handle обрабатывает PUT /courses/:id (полное обновление полей)
func (h *UpdateHandler) Handle(c *gin.Context) {
	idStr := c.Param("id")
	courseID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[ERROR] Handler.Handle: invalid ID format '%s': %v", idStr, err)
		common.HandleError(c, domain.ErrInvalidID)
		return
	}

	log.Printf("[INFO] Handler.Handle: received request to update course ID=%d", courseID)

	if !middleware.RequireAdmin(c) {
		log.Printf("[WARN] Handler.Handle: access denied for course ID=%d (not admin)", courseID)
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	var req updateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Handler.Handle: binding error for course ID=%d: %v", courseID, err)
		common.HandleError(c, err)
		return
	}

	log.Printf("[DEBUG] Handler.Handle: bound data for course ID=%d, IsActive=%v", courseID, req.IsActive)

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
		log.Printf("[ERROR] Handler.Handle: usecase execution failed for course ID=%d: %v", courseID, err)
		common.HandleError(c, err)
		return
	}

	log.Printf("[INFO] Handler.Handle: successfully updated course ID=%d", courseID)
	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

// HandleStatusPatch обрабатывает PATCH /courses/:id/status
func (h *UpdateHandler) HandleStatusPatch(c *gin.Context) {
	idStr := c.Param("id")
	courseID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[ERROR] Handler.HandleStatusPatch: invalid ID format '%s': %v", idStr, err)
		common.HandleError(c, domain.ErrInvalidID)
		return
	}

	log.Printf("[INFO] Handler.HandleStatusPatch: received request for course ID=%d", courseID)

	if !middleware.RequireAdmin(c) {
		log.Printf("[WARN] Handler.HandleStatusPatch: access denied for course ID=%d", courseID)
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Handler.HandleStatusPatch: binding error for course ID=%d: %v", courseID, err)
		common.HandleError(c, err)
		return
	}

	log.Printf("[DEBUG] Handler.HandleStatusPatch: setting course ID=%d IsActive=%v", courseID, req.IsActive)

	if err := h.usecase.UpdateStatus(c.Request.Context(), courseID, req.IsActive); err != nil {
		log.Printf("[ERROR] Handler.HandleStatusPatch: usecase execution failed for course ID=%d: %v", courseID, err)
		common.HandleError(c, err)
		return
	}

	log.Printf("[INFO] Handler.HandleStatusPatch: successfully updated status for course ID=%d", courseID)
	c.Status(http.StatusOK)
}
