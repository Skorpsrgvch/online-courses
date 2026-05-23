package lesson

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	getUC "github.com/Skorpsrgvch/online-courses/internal/usecase/lesson/get"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type lessonDTO struct {
	ID           int     `json:"id"`
	ModuleID     int     `json:"module_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key,omitempty"`
	Order        int     `json:"order"`
}

type GetHandler struct {
	usecase *getUC.Usecase
}

func NewGetHandler(usecase *getUC.Usecase) *GetHandler {
	return &GetHandler{usecase: usecase}
}

func (h *GetHandler) Handle(c *gin.Context) {
	idStr := c.Param("id")
	moduleID, err := strconv.Atoi(idStr)
	if err != nil {
		zap.L().Debug("Invalid module ID format", zap.String("id", idStr), zap.Error(err))
		common.HandleError(c, common.HttpError("invalid module ID", http.StatusBadRequest))
		return
	}

	zap.L().Debug("Getting lessons for module", zap.Int("moduleID", moduleID))

	input := getUC.Input{ModuleID: moduleID}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Failed to get lessons", zap.Int("moduleID", moduleID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	// Маппинг в DTO (можно оптимизировать, если domain.Lesson полностью совместим)
	lessons := make([]lessonDTO, 0, len(output.Lessons))
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

	zap.L().Info("Lessons retrieved", zap.Int("moduleID", moduleID), zap.Int("count", len(lessons)))
	c.JSON(http.StatusOK, gin.H{"lessons": lessons})
}
