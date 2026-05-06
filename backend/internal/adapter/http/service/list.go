package service

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/usecase/service/list"
	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	usecase *list.Usecase
}

func NewListHandler(usecase *list.Usecase) *ListHandler {
	return &ListHandler{usecase: usecase}
}

func (h *ListHandler) Handle(c *gin.Context) {
	services, err := h.usecase.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}
