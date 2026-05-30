package support

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	contactUC "github.com/Skorpsrgvch/online-courses/internal/usecase/support/contact"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ContactHandler struct {
	usecase *contactUC.Usecase
}

func NewContactHandler(usecase *contactUC.Usecase) *ContactHandler {
	return &ContactHandler{usecase: usecase}
}

func (h *ContactHandler) Handle(c *gin.Context) {
	var input contactUC.Input
	if err := c.ShouldBindJSON(&input); err != nil {
		common.HandleError(c, common.HttpError("invalid request body", http.StatusBadRequest))
		return
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Warn("Support request failed", zap.Error(err))
		common.HandleError(c, common.HttpError(err.Error(), http.StatusBadRequest))
		return
	}

	c.JSON(http.StatusOK, output)
}
