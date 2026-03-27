package module

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/module/get"
	"github.com/gin-gonic/gin"
)

type getModulesResponse struct {
	Modules []moduleDTO `json:"modules"`
}

type moduleDTO struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Order int    `json:"order"`
}

type GetHandler struct {
	usecase *get.Usecase
}

func NewGetHandler(usecase *get.Usecase) *GetHandler {
	return &GetHandler{usecase: usecase}
}

func (h *GetHandler) Handle(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, common.HttpError("invalid course ID", http.StatusBadRequest))
		return
	}

	input := get.Input{CourseID: courseID}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	var modules []moduleDTO
	for _, m := range output.Modules {
		modules = append(modules, moduleDTO{
			ID:    m.ID,
			Title: m.Title,
			Order: m.Order,
		})
	}

	c.JSON(http.StatusOK, getModulesResponse{Modules: modules})
}
