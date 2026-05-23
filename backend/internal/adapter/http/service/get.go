package service

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	getUC "github.com/Skorpsrgvch/online-courses/internal/usecase/service/get"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetByIDHandler struct {
	usecase *getUC.Usecase
}

func NewGetByIDHandler(usecase *getUC.Usecase) *GetByIDHandler {
	return &GetByIDHandler{usecase: usecase}
}

func (h *GetByIDHandler) Handle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		zap.L().Debug("Invalid service ID format", zap.String("id", idStr), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid service ID", http.StatusBadRequest))
		return
	}

	zap.L().Debug("Fetching service by ID", zap.Int("id", id))

	input := getUC.Input{ID: id}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Failed to get service", zap.Int("id", id), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
