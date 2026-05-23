package review

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	listUC "github.com/Skorpsrgvch/online-courses/internal/usecase/review/list"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ListHandler struct {
	usecase *listUC.Usecase
}

func NewListHandler(usecase *listUC.Usecase) *ListHandler {
	return &ListHandler{usecase: usecase}
}

func (h *ListHandler) Handle(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		zap.L().Debug("Invalid course ID format", zap.String("id", courseIDStr), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid course ID", http.StatusBadRequest))
		return
	}

	userID := middleware.GetUserID(c)

	zap.L().Debug("Listing reviews", zap.Int("courseID", courseID), zap.Int("userID", userID))

	input := listUC.Input{
		CourseID: courseID,
		UserID:   userID,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Failed to list reviews", zap.Int("courseID", courseID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
