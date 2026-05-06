package service

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/service/create"
	"github.com/gin-gonic/gin"
)

type CreateHandler struct {
	usecase *create.Usecase
}

func NewCreateHandler(usecase *create.Usecase) *CreateHandler {
	return &CreateHandler{usecase: usecase}
}

func (h *CreateHandler) Handle(c *gin.Context) {
	var input create.Input
	if err := c.ShouldBindJSON(&input); err != nil {
		common.HandleError(c, common.HttpError("invalid request body", http.StatusBadRequest))
		return
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}
