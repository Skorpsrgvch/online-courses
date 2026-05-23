package review

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	createUC "github.com/Skorpsrgvch/online-courses/internal/usecase/review/create"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createReviewRequest struct {
	CourseID int    `json:"course_id" binding:"required"`
	Text     string `json:"text" binding:"required,min=10"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
}

type CreateHandler struct {
	usecase *createUC.Usecase
}

func NewCreateHandler(usecase *createUC.Usecase) *CreateHandler {
	return &CreateHandler{usecase: usecase}
}

func (h *CreateHandler) Handle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, domain.ErrUnauthorized)
		return
	}

	var req createReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Debug("Invalid JSON in create review request", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Debug("Creating review", zap.Int("userID", userID), zap.Int("courseID", req.CourseID), zap.Int("rating", req.Rating))

	input := createUC.Input{
		UserID:   userID,
		CourseID: req.CourseID,
		Text:     req.Text,
		Rating:   req.Rating,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Info("Failed to create review", zap.Int("userID", userID), zap.Int("courseID", req.CourseID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Review created successfully", zap.Int("userID", userID), zap.Int("courseID", req.CourseID))
	c.Status(http.StatusCreated)
}
