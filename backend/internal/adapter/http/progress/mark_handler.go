package progress

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/progress/mark"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MarkHandler struct {
	usecase *mark.Usecase
}

func NewMarkHandler(usecase *mark.Usecase) *MarkHandler {
	return &MarkHandler{usecase: usecase}
}

func (h *MarkHandler) Handle(c *gin.Context) {
	lessonID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, common.HttpError("invalid lesson ID", http.StatusBadRequest))
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, domain.ErrUnauthorized)
		return
	}

	zap.L().Debug("Marking lesson as completed",
		zap.Int("user_id", userID),
		zap.Int("lesson_id", lessonID))

	input := mark.Input{UserID: userID, LessonID: lessonID}
	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Warn("Failed to mark lesson completed", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Lesson marked as completed", zap.Int("lesson_id", lessonID))
	c.Status(http.StatusOK)
}
