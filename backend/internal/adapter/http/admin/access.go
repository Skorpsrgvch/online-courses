package admin

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/admin/access"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GrantAccessHandler struct {
	usecase *access.Usecase
}

func NewGrantAccessHandler(usecase *access.Usecase) *GrantAccessHandler {
	return &GrantAccessHandler{usecase: usecase}
}

func (h *GrantAccessHandler) Handle(c *gin.Context) {
	var input access.Input

	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Debug("Invalid JSON in grant access", zap.Error(err))
		common.HandleError(c, common.HttpError("Неверный формат данных", http.StatusBadRequest))
		return
	}

	if input.UserID <= 0 || input.CourseID <= 0 {
		common.HandleError(c, common.HttpError("ID пользователя и курса должны быть положительными", http.StatusBadRequest))
		return
	}

	zap.L().Info("Granting access request", zap.Int("userID", input.UserID), zap.Int("courseID", input.CourseID))

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Grant access failed", zap.Int("userID", input.UserID), zap.Int("courseID", input.CourseID), zap.Error(err))

		// Маппинг ошибок на статусы
		status := http.StatusInternalServerError
		if err.Error() == "пользователь не найден" || err.Error() == "курс не найден" {
			status = http.StatusNotFound
		}
		common.HandleError(c, common.HttpError(err.Error(), status))
		return
	}

	if !output.Success {
		zap.L().Info("Access already exists", zap.Int("userID", input.UserID), zap.Int("courseID", input.CourseID))
		c.JSON(http.StatusConflict, gin.H{"error": output.Message})
		return
	}

	zap.L().Info("Access granted successfully", zap.Int("userID", input.UserID), zap.Int("courseID", input.CourseID))
	c.JSON(http.StatusOK, gin.H{
		"message":   output.Message,
		"user_id":   input.UserID,
		"course_id": input.CourseID,
	})
}
