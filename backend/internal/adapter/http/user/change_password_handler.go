package user

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/user/change_password"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ChangePasswordHandler struct {
	usecase *change_password.Usecase
}

func NewChangePasswordHandler(usecase *change_password.Usecase) *ChangePasswordHandler {
	return &ChangePasswordHandler{usecase: usecase}
}

func (h *ChangePasswordHandler) Handle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, common.HttpError("требуется авторизация", http.StatusUnauthorized))
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	zap.L().Debug("Change password request", zap.Int("user_id", userID))

	input := change_password.Input{
		UserID:      userID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Warn("Change password failed", zap.Int("user_id", userID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Password changed successfully", zap.Int("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": output.Message})
}
