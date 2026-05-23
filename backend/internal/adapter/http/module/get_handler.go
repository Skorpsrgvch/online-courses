package module

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	getUC "github.com/Skorpsrgvch/online-courses/internal/usecase/module/get"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type moduleDTO struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Order int    `json:"order"`
}

type GetHandler struct {
	usecase *getUC.Usecase
}

func NewGetHandler(usecase *getUC.Usecase) *GetHandler {
	return &GetHandler{usecase: usecase}
}

func (h *GetHandler) Handle(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zap.L().Debug("Invalid course ID format", zap.String("id", c.Param("id")), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid course ID", http.StatusBadRequest))
		return
	}

	zap.L().Debug("Getting modules for course", zap.Int("courseID", courseID))

	input := getUC.Input{CourseID: courseID}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Failed to get modules", zap.Int("courseID", courseID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	modules := make([]moduleDTO, 0, len(output.Modules))
	for _, m := range output.Modules {
		modules = append(modules, moduleDTO{
			ID:    m.ID,
			Title: m.Title,
			Order: m.Order,
		})
	}

	c.JSON(http.StatusOK, gin.H{"modules": modules})
}
