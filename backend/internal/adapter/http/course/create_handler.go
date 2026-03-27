package course

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
	"github.com/gin-gonic/gin"
)

type createCourseRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	Price       int    `json:"price"`
}

type CreateHandler struct {
	usecase *create.Usecase
}

func NewCreateHandler(usecase *create.Usecase) *CreateHandler {
	return &CreateHandler{usecase: usecase}
}

func (h *CreateHandler) Handle(c *gin.Context) {
	var req createCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	// ← Используем функции из middleware
	userID := middleware.GetUserID(c)
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	input := create.Input{
		Title:       req.Title,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		Price:       req.Price,
		AuthorID:    userID,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}
