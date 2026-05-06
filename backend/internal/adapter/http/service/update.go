package service

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/service/update"
	"github.com/gin-gonic/gin"
)

type UpdateHandler struct {
	usecase *update.Usecase
}

func NewUpdateHandler(usecase *update.Usecase) *UpdateHandler {
	return &UpdateHandler{usecase: usecase}
}

func (h *UpdateHandler) Handle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.HandleError(c, common.HttpError("invalid service ID", http.StatusBadRequest))
		return
	}

	var input update.Input
	if err := c.ShouldBindJSON(&input); err != nil {
		common.HandleError(c, common.HttpError("invalid request body", http.StatusBadRequest))
		return
	}
	input.ID = id

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
