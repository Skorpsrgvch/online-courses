package user

import (
	"net/http"
	"strings"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/user/update_profile"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UpdateProfileRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

type UpdateProfileHandler struct {
	usecase *update_profile.Usecase
}

func NewUpdateProfileHandler(usecase *update_profile.Usecase) *UpdateProfileHandler {
	return &UpdateProfileHandler{usecase: usecase}
}

func (h *UpdateProfileHandler) Handle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, common.HttpError("требуется авторизация", http.StatusUnauthorized))
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, common.HttpError("неверный формат запроса", http.StatusBadRequest))
		return
	}

	if req.Email != nil && strings.TrimSpace(*req.Email) == "" {
		common.HandleError(c, common.HttpError("email не может быть пустым", http.StatusBadRequest))
		return
	}

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		common.HandleError(c, common.HttpError("имя не может быть пустым", http.StatusBadRequest))
		return
	}

	if req.Name == nil && req.Email == nil {
		common.HandleError(c, common.HttpError("укажите имя или email для обновления", http.StatusBadRequest))
		return
	}

	zap.L().Debug("Update profile request",
		zap.Int("user_id", userID),
		zap.Bool("name_changed", req.Name != nil),
		zap.Bool("email_changed", req.Email != nil))

	input := update_profile.Input{
		UserID:   userID,
		NewName:  req.Name,
		NewEmail: req.Email,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Warn("Update profile failed", zap.Int("user_id", userID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Profile updated", zap.Int("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": output.Message})
}
