package course

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	getFullUC "github.com/Skorpsrgvch/online-courses/internal/usecase/course/getfull"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetFullHandler struct {
	usecase *getFullUC.Usecase
}

func NewGetFullHandler(usecase *getFullUC.Usecase) *GetFullHandler {
	return &GetFullHandler{usecase: usecase}
}

func (h *GetFullHandler) Handle(c *gin.Context) {
	idStr := c.Param("id")
	courseID, err := strconv.Atoi(idStr)
	if err != nil {
		common.HandleError(c, common.HttpError("Неверный ID курса", http.StatusBadRequest))
		return
	}

	userID := 0
	role := ""

	if uidRaw, exists := c.Get("user_id"); exists {
		if id, ok := uidRaw.(int); ok {
			userID = id
		}
	}

	if r, exists := c.Get("role"); exists && r != nil {
		if roleVal, ok := r.(string); ok {
			role = roleVal
		}
	}

	zap.L().Debug("Get full course request", zap.Int("courseID", courseID), zap.Int("userID", userID), zap.String("role", role))

	input := getFullUC.Input{
		CourseID: courseID,
		UserID:   userID,
		Role:     role,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Get full course failed", zap.Int("courseID", courseID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	if output.Course != nil {
		output.Course.IsPurchased = output.IsPurchased
	}

	zap.L().Info("Get full course successful", zap.Int("courseID", courseID), zap.Bool("isPurchased", output.IsPurchased))
	c.JSON(http.StatusOK, output)
}
