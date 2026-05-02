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

type updateLessonRequest struct {
	ID           int     `json:"id"`
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key"`
	Order        int     `json:"order"`
}

type updateModuleRequest struct {
	ID      int                   `json:"id"`
	Title   string                `json:"title" binding:"required"`
	Order   int                   `json:"order"`
	Lessons []updateLessonRequest `json:"lessons"`
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
	usecase *updatefullcourse.Usecase
}

func NewUpdateFullCourseHandler(usecase *updatefullcourse.Usecase) *UpdateFullCourseHandler {
	return &UpdateFullCourseHandler{usecase: usecase}
}

func (h *UpdateFullCourseHandler) Handle(c *gin.Context) {
	// 1. Получаем ID курса из URL параметра
	idStr := c.Param("id")
	courseID, err := strconv.Atoi(idStr)
	if err != nil {
		common.HandleError(c, domain.ErrInvalidID)
		return
	}

	var req updateFullCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	userID := middleware.GetUserID(c)

	// Конвертируем запрос в формат юзкейса
	modules := make([]updatefullcourse.ModuleInput, len(req.Modules))
	for i, mod := range req.Modules {
		lessons := make([]updatefullcourse.LessonInput, len(mod.Lessons))
		for j, l := range mod.Lessons {
			lessons[j] = updatefullcourse.LessonInput{
				ID:           l.ID,
				Title:        l.Title,
				Description:  l.Description,
				VideoEmbedID: l.VideoEmbedID,
				PrivateKey:   l.PrivateKey,
				Order:        l.Order,
			}
		}
		modules[i] = updatefullcourse.ModuleInput{
			ID:      mod.ID,
			Title:   mod.Title,
			Order:   mod.Order,
			Lessons: lessons,
		}
	}

	input := updatefullcourse.Input{
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
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
