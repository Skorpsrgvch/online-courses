package payment

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/payment/create"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateHandler struct {
	usecase *create.UseCase
}

func NewCreateHandler(usecase *create.UseCase) *CreateHandler {
	return &CreateHandler{usecase: usecase}
}

func (h *CreateHandler) Handle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, common.HttpError("unauthorized", http.StatusUnauthorized))
		return
	}

	var input create.Input
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Debug("Invalid JSON in create payment", zap.Error(err))
		common.HandleError(c, common.HttpError("invalid request body", http.StatusBadRequest))
		return
	}

	input.UserID = userID

	zap.L().Info("Creating payment",
		zap.Int("user_id", userID),
		zap.Int("course_id", input.CourseID))

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Warn("Create payment failed", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, output)
}
