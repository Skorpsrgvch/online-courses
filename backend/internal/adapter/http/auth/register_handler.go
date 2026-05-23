package auth

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	registerUC "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/register"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterHandler struct {
	usecase *registerUC.Usecase
}

func NewRegisterHandler(usecase *registerUC.Usecase) *RegisterHandler {
	return &RegisterHandler{usecase: usecase}
}

func (h *RegisterHandler) Handle(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	zap.L().Debug("Registration request", zap.String("email", req.Email))

	input := registerUC.Input{
		Email:    req.Email,
		FullName: req.FullName,
		Password: req.Password,
		Role:     "user",
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Info("Registration failed", zap.String("email", req.Email), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("User registered successfully", zap.String("email", req.Email))
	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}
