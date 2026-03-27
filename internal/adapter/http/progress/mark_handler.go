package progress

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/progress/mark"
	"github.com/gin-gonic/gin"
)

type MarkHandler struct {
	usecase *mark.Usecase
}

func NewMarkHandler(usecase *mark.Usecase) *MarkHandler {
	return &MarkHandler{usecase: usecase}
}

func (h *MarkHandler) Handle(c *gin.Context) {
	lessonID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, common.HttpError("invalid lesson ID", http.StatusBadRequest))
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, domain.ErrUnauthorized)
		return
	}

	input := mark.Input{
		UserID:   userID,
		LessonID: lessonID,
	}

	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
