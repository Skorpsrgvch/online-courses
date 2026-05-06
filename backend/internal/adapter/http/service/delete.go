package service

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/service/delete"
	"github.com/gin-gonic/gin"
)

type DeleteHandler struct {
	usecase *delete.Usecase
}

func NewDeleteHandler(usecase *delete.Usecase) *DeleteHandler {
	return &DeleteHandler{usecase: usecase}
}

func (h *DeleteHandler) Handle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.HandleError(c, common.HttpError("invalid service ID", http.StatusBadRequest))
		return
	}

	input := delete.Input{ID: id}
	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
