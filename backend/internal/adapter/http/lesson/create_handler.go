package lesson

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	createUC "github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/create"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createLessonRequest struct {
	ModuleID     int     `json:"module_id" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key"`
	Order        int     `json:"order"`
}

type CreateHandler struct {
	usecase *createUC.Usecase
}

func NewCreateHandler(usecase *createUC.Usecase) *CreateHandler {
	return &CreateHandler{usecase: usecase}
}

func (h *CreateHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, common.HttpError("admin access required", http.StatusForbidden))
		return
	}

	var req createLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Debug("Invalid JSON in create lesson request", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Creating new lesson", zap.Int("moduleID", req.ModuleID), zap.String("title", req.Title))

	input := createUC.Input{
		ModuleID:     req.ModuleID,
		Title:        req.Title,
		Description:  req.Description,
		VideoEmbedID: req.VideoEmbedID,
		PrivateKey:   req.PrivateKey,
		Order:        req.Order,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Failed to create lesson", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Lesson created successfully", zap.Int("moduleID", req.ModuleID))
	c.Status(http.StatusCreated)
}
