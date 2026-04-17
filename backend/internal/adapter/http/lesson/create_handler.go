package lesson

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/create"
	"github.com/gin-gonic/gin"
)

type createLessonRequest struct {
	ModuleID       int    `json:"module_id" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description"`
	LessonType     string `json:"lesson_type" binding:"required,oneof=video article"`
	VideoEmbedID   string `json:"video_embed_id"`
	ArticleContent string `json:"article_content"`
	Order          int    `json:"order"`
}

type CreateHandler struct {
	usecase *create.Usecase
}

func NewCreateHandler(usecase *create.Usecase) *CreateHandler {
	return &CreateHandler{usecase: usecase}
}

func (h *CreateHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	var req createLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	input := create.Input{
		ModuleID:       req.ModuleID,
		Title:          req.Title,
		Description:    req.Description,
		LessonType:     domain.LessonType(req.LessonType),
		VideoEmbedID:   req.VideoEmbedID,
		ArticleContent: req.ArticleContent,
		Order:          req.Order,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}
