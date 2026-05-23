package review

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	adminDeleteUC "github.com/Skorpsrgvch/online-courses/internal/usecase/review/admin_delete"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminDeleteHandler struct {
	usecase *adminDeleteUC.Usecase
}

func NewAdminDeleteHandler(usecase *adminDeleteUC.Usecase) *AdminDeleteHandler {
	return &AdminDeleteHandler{usecase: usecase}
}

func (h *AdminDeleteHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, common.HttpError("admin access required", http.StatusForbidden))
		return
	}

	reviewIDStr := c.Param("id")
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		zap.L().Debug("Invalid review ID format", zap.String("id", reviewIDStr), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid review ID", http.StatusBadRequest))
		return
	}

	zap.L().Debug("Admin deleting review", zap.Int("reviewID", reviewID))

	input := adminDeleteUC.Input{ReviewID: reviewID}
	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Failed to delete review", zap.Int("reviewID", reviewID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Review deleted by admin", zap.Int("reviewID", reviewID))
	c.Status(http.StatusNoContent)
}
