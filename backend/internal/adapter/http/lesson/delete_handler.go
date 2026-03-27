package lesson

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/delete"
	"github.com/gin-gonic/gin"
)

type DeleteHandler struct {
	usecase *delete.Usecase
}

func NewDeleteHandler(usecase *delete.Usecase) *DeleteHandler {
	return &DeleteHandler{usecase: usecase}
}

func (h *DeleteHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, common.HttpError("admin access required", http.StatusForbidden))
		return
	}

	lessonID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, common.HttpError("invalid lesson ID", http.StatusBadRequest))
		return
	}

	input := delete.Input{ID: lessonID}
	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
