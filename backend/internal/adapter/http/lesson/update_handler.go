package lesson

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	updateUC "github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/update"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type updateLessonRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key"`
	Order        int     `json:"order"`
}

type UpdateHandler struct {
	usecase *updateUC.Usecase
}

func NewUpdateHandler(usecase *updateUC.Usecase) *UpdateHandler {
	return &UpdateHandler{usecase: usecase}
}

func (h *UpdateHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, common.HttpError("admin access required", http.StatusForbidden))
		return
	}

	idStr := c.Param("id")
	lessonID, err := strconv.Atoi(idStr)
	if err != nil {
		zap.L().Debug("Invalid lesson ID format", zap.String("id", idStr), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid lesson ID", http.StatusBadRequest))
		return
	}

	var req updateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Debug("Invalid JSON in update lesson request", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Updating lesson", zap.Int("lessonID", lessonID))

	input := updateUC.Input{
		ID:           lessonID,
		Title:        req.Title,
		Description:  req.Description,
		VideoEmbedID: req.VideoEmbedID,
		PrivateKey:   req.PrivateKey,
		Order:        req.Order,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Failed to update lesson", zap.Int("lessonID", lessonID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Lesson updated successfully", zap.Int("lessonID", lessonID))
	c.Status(http.StatusOK)
}
