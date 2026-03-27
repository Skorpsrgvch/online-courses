package auth

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/register"
	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterHandler struct {
	usecase *register.Usecase
}

func NewRegisterHandler(usecase *register.Usecase) *RegisterHandler {
	return &RegisterHandler{usecase: usecase}
}

func (h *RegisterHandler) Handle(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleError(c, err)
		return
	}

	input := register.Input{
		Email:    req.Email,
		FullName: req.FullName,
		Password: req.Password,
		Role:     "user", // только user может регистрироваться
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}
