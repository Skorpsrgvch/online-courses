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

type updateModuleRequest struct {
	ID      int                   `json:"id"`
	Title   string                `json:"title" binding:"required"`
	Order   int                   `json:"order"`
	Lessons []updateLessonRequest `json:"lessons"`
}

type updateLessonRequest struct {
	ID           int     `json:"id"`
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key"`
	Order        int     `json:"order"`
}

type updateFullCourseRequest struct {
	Title             string                `json:"title" binding:"required"`
	Description       string                `json:"description"`
	IsPublic          bool                  `json:"is_public"`
	Price             int                   `json:"price"`
	IsActive          bool                  `json:"is_active"`
	CoverImageURL     string                `json:"cover_image_url"`
	Contraindications string                `json:"contraindications"`
	Recommendations   string                `json:"recommendations"`
	TargetAudience    string                `json:"target_audience"`
	CourseBasis       string                `json:"course_basis"`
	ClassBasis        string                `json:"class_basis"`
	Bonuses           []domain.BonusItem    `json:"bonuses"`
	Modules           []updateModuleRequest `json:"modules"`
}

type UpdateFullCourseHandler struct {
	usecase *updateFullUC.Usecase
}

func NewUpdateFullCourseHandler(usecase *updateFullUC.Usecase) *UpdateFullCourseHandler {
	return &UpdateFullCourseHandler{usecase: usecase}
}

func (h *UpdateFullCourseHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, domain.ErrInvalidID)
		return
	}

	var req updateFullCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	userID := middleware.GetUserID(c)
	zap.L().Info("Update full course request", zap.Int("courseID", courseID), zap.Int("authorID", userID))

	modules := make([]updateFullUC.ModuleInput, len(req.Modules))
	for i, mod := range req.Modules {
		lessons := make([]updateFullUC.LessonInput, len(mod.Lessons))
		for j, l := range mod.Lessons {
			lessons[j] = updateFullUC.LessonInput{
				ID:           l.ID,
				Title:        l.Title,
				Description:  l.Description,
				VideoEmbedID: l.VideoEmbedID,
				PrivateKey:   l.PrivateKey,
				Order:        l.Order,
			}
		}
		modules[i] = updateFullUC.ModuleInput{
			ID:      mod.ID,
			Title:   mod.Title,
			Order:   mod.Order,
			Lessons: lessons,
		}
	}

	input := updateFullUC.Input{
		CourseID:          courseID,
		Title:             req.Title,
		Description:       req.Description,
		IsPublic:          req.IsPublic,
		Price:             req.Price,
		IsActive:          req.IsActive,
		AuthorID:          userID,
		CoverImageURL:     req.CoverImageURL,
		Contraindications: req.Contraindications,
		Recommendations:   req.Recommendations,
		TargetAudience:    req.TargetAudience,
		CourseBasis:       req.CourseBasis,
		ClassBasis:        req.ClassBasis,
		Bonuses:           req.Bonuses,
		Modules:           modules,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Update full course failed", zap.Int("courseID", courseID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Full course updated successfully", zap.Int("courseID", courseID))
	c.Status(http.StatusOK)
}
