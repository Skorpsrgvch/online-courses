package course

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	enrollFreeUC "github.com/Skorpsrgvch/online-courses/internal/usecase/course/enroll_free"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EnrollFreeHandler struct {
	usecase *enrollFreeUC.Usecase
}

func NewEnrollFreeHandler(usecase *enrollFreeUC.Usecase) *EnrollFreeHandler {
	return &EnrollFreeHandler{usecase: usecase}
}

func (h *EnrollFreeHandler) Handle(c *gin.Context) {
	rawID := c.Param("id")
	courseID, err := strconv.Atoi(rawID)
	if err != nil {
		zap.L().Debug("Invalid course ID format", zap.String("rawID", rawID), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid course ID", http.StatusBadRequest))
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, common.HttpError("unauthorized", http.StatusUnauthorized))
		return
	}

	var req struct {
		Price int `json:"price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Debug("Binding JSON failed for enroll free", zap.Error(err))
		common.HandleError(c, common.HttpError("price is required (json)", http.StatusBadRequest))
		return
	}

	zap.L().Info("Enroll free request", zap.Int("userID", userID), zap.Int("courseID", courseID), zap.Int("price", req.Price))

	input := enrollFreeUC.Input{
		UserID:      userID,
		CourseID:    courseID,
		CoursePrice: req.Price,
	}

	if _, err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Warn("Enroll free failed", zap.Int("userID", userID), zap.Int("courseID", courseID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Enroll free successful", zap.Int("userID", userID), zap.Int("courseID", courseID))
	c.JSON(http.StatusOK, gin.H{"message": "Доступ предоставлен"})
}
