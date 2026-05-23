package service

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	createUC "github.com/Skorpsrgvch/online-courses/internal/usecase/service/create"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateHandler struct {
	usecase *createUC.Usecase
}

func NewCreateHandler(usecase *createUC.Usecase) *CreateHandler {
	return &CreateHandler{usecase: usecase}
}

func (h *CreateHandler) Handle(c *gin.Context) {
	var input createUC.Input
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Debug("Invalid JSON in create service request", zap.Error(err))
		common.HandleError(c, common.HttpError("invalid request body", http.StatusBadRequest))
		return
	}

	zap.L().Info("Creating new service", zap.String("title", input.Title))

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		zap.L().Error("Failed to create service", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	zap.L().Info("Service created successfully")
	c.Status(http.StatusCreated)
}
