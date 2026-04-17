package course

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
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
