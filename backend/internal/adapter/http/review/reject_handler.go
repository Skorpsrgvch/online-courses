package review

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	rejectUC "github.com/Skorpsrgvch/online-courses/internal/usecase/review/reject"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type rejectReviewRequest struct {
	Reason string `json:"reason" binding:"required,min=5"`
}

type RejectHandler struct {
	usecase *rejectUC.Usecase
}

func NewRejectHandler(usecase *rejectUC.Usecase) *RejectHandler {
	return &RejectHandler{usecase: usecase}
}

func (h *RejectHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	reviewIDStr := c.Param("id")
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		zap.L().Debug("Invalid review ID format", zap.String("id", reviewIDStr), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid review ID", http.StatusBadRequest))
		return
	}

	var req rejectReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Debug("Invalid JSON in reject review request", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Debug("Rejecting review", zap.Int("reviewID", reviewID), zap.String("reason", req.Reason))

	input := rejectUC.Input{
		ReviewID: reviewID,
		Reason:   req.Reason,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Failed to reject review", zap.Int("reviewID", reviewID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Review rejected", zap.Int("reviewID", reviewID))
	c.Status(http.StatusOK)
}
