package lesson

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/get"
	"github.com/gin-gonic/gin"
)

type getLessonsResponse struct {
	Lessons []lessonDTO `json:"lessons"`
}

type lessonDTO struct {
	ID           int     `json:"id"`
	ModuleID     int     `json:"module_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key"`
	Order        int     `json:"order"`
}

type GetHandler struct {
	usecase *get.Usecase
}

func NewGetHandler(usecase *get.Usecase) *GetHandler {
	return &GetHandler{usecase: usecase}
}

func (h *GetHandler) Handle(c *gin.Context) {
	moduleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, common.HttpError("invalid module ID", http.StatusBadRequest))
		return
	}

	input := get.Input{ModuleID: moduleID}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	var lessons []lessonDTO
	for _, l := range output.Lessons {
		lessons = append(lessons, lessonDTO{
			ID:           l.ID,
			ModuleID:     l.ModuleID,
			Title:        l.Title,
			Description:  l.Description,
			VideoEmbedID: l.VideoEmbedID,
			PrivateKey:   l.PrivateKey,
			Order:        l.Order,
		})
	}

	c.JSON(http.StatusOK, getLessonsResponse{Lessons: lessons})
}
