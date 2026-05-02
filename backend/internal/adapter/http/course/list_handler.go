package course

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/list"
	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	usecase *list.Usecase
}

func NewListHandler(usecase *list.Usecase) *ListHandler {
	return &ListHandler{usecase: usecase}
}

func (h *ListHandler) Handle(c *gin.Context) {
	output, err := h.usecase.Execute(c.Request.Context())
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output.Courses)
}

func (h *ListHandler) HandleAdmin(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	courses, err := h.usecase.ExecuteAdmin(c.Request.Context())
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}
