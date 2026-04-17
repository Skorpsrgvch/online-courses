package review

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/review/pending"
	"github.com/gin-gonic/gin"
)

type PendingHandler struct {
	usecase *pending.Usecase
}

func NewPendingHandler(usecase *pending.Usecase) *PendingHandler {
	return &PendingHandler{usecase: usecase}
}

func (h *PendingHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	input := pending.Input{}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
