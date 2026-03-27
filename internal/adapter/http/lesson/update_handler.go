package lesson

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/update"
	"github.com/gin-gonic/gin"
)

type updateLessonRequest struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	VideoEmbedID string `json:"video_embed_id" binding:"required"`
	Order        int    `json:"order"`
}

type UpdateHandler struct {
	usecase *update.Usecase
}

func NewUpdateHandler(usecase *update.Usecase) *UpdateHandler {
	return &UpdateHandler{usecase: usecase}
}

func (h *UpdateHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, common.HttpError("admin access required", http.StatusForbidden))
		return
	}

	lessonID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, common.HttpError("invalid lesson ID", http.StatusBadRequest))
		return
	}

	var req updateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	input := update.Input{
		ID:           lessonID,
		Title:        req.Title,
		Description:  req.Description,
		VideoEmbedID: req.VideoEmbedID,
		Order:        req.Order,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
