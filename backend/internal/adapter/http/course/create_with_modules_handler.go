package course

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
	createwithmodules "github.com/Skorpsrgvch/online-courses/internal/usecase/course/createwithmodules"
	"github.com/gin-gonic/gin"
)

type moduleRequest struct {
	Title   string          `json:"title" binding:"required"`
	Order   int             `json:"order"`
	Lessons []lessonRequest `json:"lessons"`
}

type lessonRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key"`
	Order        int     `json:"order"`
}

type createCourseWithModulesRequest struct {
	Title             string             `json:"title" binding:"required"`
	Description       string             `json:"description"`
	IsPublic          bool               `json:"is_public"`
	Price             int                `json:"price"`
	CoverImageURL     string             `json:"cover_image_url"`
	Contraindications string             `json:"contraindications"`
	Recommendations   string             `json:"recommendations"`
	TargetAudience    string             `json:"target_audience"`
	CourseBasis       string             `json:"course_basis"`
	ClassBasis        string             `json:"class_basis"`
	Bonuses           []domain.BonusItem `json:"bonuses"`
	Modules           []moduleRequest    `json:"modules"`
}

type CreateWithModulesHandler struct {
	usecase *createwithmodules.Usecase
}

func NewCreateWithModulesHandler(usecase *createwithmodules.Usecase) *CreateWithModulesHandler {
	return &CreateWithModulesHandler{usecase: usecase}
}

func (h *CreateWithModulesHandler) Handle(c *gin.Context) {
	var req createCourseWithModulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	userID := middleware.GetUserID(c)

	// Конвертируем модули из запроса в domain
	modules := make([]create.ModuleInput, len(req.Modules))
	for i, mod := range req.Modules {
		lessons := make([]create.LessonInput, len(mod.Lessons))
		for j, l := range mod.Lessons {
			lessons[j] = create.LessonInput{
				Title:        l.Title,
				Description:  l.Description,
				VideoEmbedID: l.VideoEmbedID,
				PrivateKey:   l.PrivateKey,
				Order:        l.Order,
			}
		}
		modules[i] = create.ModuleInput{
			Title:   mod.Title,
			Order:   mod.Order,
			Lessons: lessons,
		}
	}

	input := createwithmodules.Input{
		Title:             req.Title,
		Description:       req.Description,
		IsPublic:          req.IsPublic,
		Price:             req.Price,
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

	c.Status(http.StatusCreated)
}
