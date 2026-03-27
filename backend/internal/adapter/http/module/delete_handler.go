package module

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/module/delete"
	"github.com/gin-gonic/gin"
)

type DeleteHandler struct {
	usecase *delete.Usecase
}

func NewDeleteHandler(usecase *delete.Usecase) *DeleteHandler {
	return &DeleteHandler{usecase: usecase}
}

func (h *DeleteHandler) Handle(c *gin.Context) {
	moduleID, _ := strconv.Atoi(c.Param("id"))

	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	input := delete.Input{ID: moduleID}
	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
