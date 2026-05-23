package review

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	deleteUC "github.com/Skorpsrgvch/online-courses/internal/usecase/review/delete"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeleteHandler struct {
	usecase *deleteUC.Usecase
}

func NewDeleteHandler(usecase *deleteUC.Usecase) *DeleteHandler {
	return &DeleteHandler{usecase: usecase}
}

func (h *DeleteHandler) Handle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, common.HttpError("unauthorized", http.StatusUnauthorized))
		return
	}

	reviewIDStr := c.Param("id")
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		zap.L().Debug("Invalid review ID format", zap.String("id", reviewIDStr), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid review ID", http.StatusBadRequest))
		return
	}

	zap.L().Debug("Deleting review", zap.Int("reviewID", reviewID), zap.Int("userID", userID))

	input := deleteUC.Input{ReviewID: reviewID}
	if _, err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Info("Failed to delete review", zap.Int("reviewID", reviewID), zap.Int("userID", userID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Review deleted", zap.Int("reviewID", reviewID))
	c.Status(http.StatusNoContent)
}
