package user

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	usercourses "github.com/Skorpsrgvch/online-courses/internal/usecase/user/courses"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CoursesHandler struct {
	usecase *usercourses.Usecase
}

func NewCoursesHandler(usecase *usercourses.Usecase) *CoursesHandler {
	return &CoursesHandler{usecase: usecase}
}

func (h *CoursesHandler) Handle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, domain.ErrUnauthorized)
		return
	}

	input := usercourses.Input{UserID: userID}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Failed to get user courses", zap.Int("user_id", userID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
