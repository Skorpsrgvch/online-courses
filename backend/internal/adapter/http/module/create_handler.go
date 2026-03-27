package module

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/module/create"
)

type createModuleRequest struct {
	CourseID int    `json:"course_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Order    int    `json:"order"`
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

	var req createModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	input := create.Input{
		CourseID: req.CourseID,
		Title:    req.Title,
		Order:    req.Order,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}
