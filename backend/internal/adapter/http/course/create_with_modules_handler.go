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
	Title   string           `json:"title" binding:"required"`
	Order   int              `json:"order"`
	Lessons []lessonRequest  `json:"lessons"`
}

type lessonRequest struct {
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description"`
	LessonType     string `json:"lesson_type" binding:"required,oneof=video article"`
	VideoEmbedID   string `json:"video_embed_id"`
	ArticleContent string `json:"article_content"`
	Order          int    `json:"order"`
}

type createCourseWithModulesRequest struct {
	Title         string          `json:"title" binding:"required"`
	Description   string          `json:"description"`
	IsPublic      bool            `json:"is_public"`
	Price         int             `json:"price"`
	CoverImageURL string          `json:"cover_image_url"`
	Modules       []moduleRequest `json:"modules"`
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

	userID := middleware.GetUserID(c)
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	// Конвертируем модули из запроса в domain
	modules := make([]create.ModuleInput, len(req.Modules))
	for i, mod := range req.Modules {
		lessons := make([]create.LessonInput, len(mod.Lessons))
		for j, l := range mod.Lessons {
			lessons[j] = create.LessonInput{
				Title:          l.Title,
				Description:    l.Description,
				LessonType:     domain.LessonType(l.LessonType),
				VideoEmbedID:   l.VideoEmbedID,
				ArticleContent: l.ArticleContent,
				Order:          l.Order,
			}
		}
		modules[i] = create.ModuleInput{
			Title:   mod.Title,
			Order:   mod.Order,
			Lessons: lessons,
		}
	}

	input := createwithmodules.Input{
		Title:         req.Title,
		Description:   req.Description,
		IsPublic:      req.IsPublic,
		Price:         req.Price,
		AuthorID:      userID,
		CoverImageURL: req.CoverImageURL,
		Modules:       modules,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}
