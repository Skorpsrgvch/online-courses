package course

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/get"
	"github.com/gin-gonic/gin"
)

type GetHandler struct {
	usecase *get.Usecase
}

func NewGetHandler(usecase *get.Usecase) *GetHandler {
	return &GetHandler{usecase: usecase}
}

func (h *GetHandler) Handle(c *gin.Context) {
	courseID, _ := strconv.Atoi(c.Param("id"))
	userID := 0
	role := ""
	if uid := middleware.GetUserID(c); uid != 0 {
		userID = uid
		if r, ok := c.Get("role"); ok {
			role = r.(string)
		}
	}

	input := get.Input{
		CourseID: courseID,
		UserID:   userID,
		Role:     role,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output.Course)
}
