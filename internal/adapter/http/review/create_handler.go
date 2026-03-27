package review

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/review/create"
	"github.com/gin-gonic/gin"
)

type createReviewRequest struct {
	CourseID int    `json:"course_id" binding:"required"`
	Text     string `json:"text" binding:"required,min=10"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
}

type CreateHandler struct {
	usecase *create.Usecase
}

func NewCreateHandler(usecase *create.Usecase) *CreateHandler {
	return &CreateHandler{usecase: usecase}
}

func (h *CreateHandler) Handle(c *gin.Context) {
	var req createReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, domain.ErrUnauthorized)
		return
	}

	input := create.Input{
		UserID:   userID,
		CourseID: req.CourseID,
		Text:     req.Text,
		Rating:   req.Rating,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}
