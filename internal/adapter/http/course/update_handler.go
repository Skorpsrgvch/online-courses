package course

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/update"
	"github.com/gin-gonic/gin"
)

type updateCourseRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	Price       int    `json:"price"`
	IsActive    bool   `json:"is_active"`
}

type UpdateHandler struct {
	usecase *update.Usecase
}

func NewUpdateHandler(usecase *update.Usecase) *UpdateHandler {
	return &UpdateHandler{usecase: usecase}
}

func (h *UpdateHandler) Handle(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, common.HttpError("invalid course ID", http.StatusBadRequest))
		return
	}

	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	var req updateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	userID := middleware.GetUserID(c)
	input := update.Input{
		ID:          courseID,
		Title:       req.Title,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		Price:       req.Price,
		IsActive:    req.IsActive,
		AuthorID:    userID,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
