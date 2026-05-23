package review

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	approveUC "github.com/Skorpsrgvch/online-courses/internal/usecase/review/approve"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApproveHandler struct {
	usecase *approveUC.Usecase
}

func NewApproveHandler(usecase *approveUC.Usecase) *ApproveHandler {
	return &ApproveHandler{usecase: usecase}
}

func (h *ApproveHandler) Handle(c *gin.Context) {
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

	zap.L().Debug("Approving review", zap.Int("reviewID", reviewID))

	input := approveUC.Input{ReviewID: reviewID}
	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Failed to approve review", zap.Int("reviewID", reviewID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Review approved", zap.Int("reviewID", reviewID))
	c.Status(http.StatusOK)
}
