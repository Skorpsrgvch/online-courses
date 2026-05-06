package service

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/service/get"
	"github.com/gin-gonic/gin"
)

type GetByIDHandler struct {
	usecase *get.Usecase
}

func NewGetByIDHandler(usecase *get.Usecase) *GetByIDHandler {
	return &GetByIDHandler{usecase: usecase}
}

func (h *GetByIDHandler) Handle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.HandleError(c, common.HttpError("invalid service ID", http.StatusBadRequest))
		return
	}

	input := get.Input{
		ID: id,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
