package course

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/getfull"
	"github.com/gin-gonic/gin"
)

type GetFullHandler struct {
	usecase *getfull.Usecase
}

func NewGetFullHandler(usecase *getfull.Usecase) *GetFullHandler {
	return &GetFullHandler{usecase: usecase}
}

func (h *GetFullHandler) Handle(c *gin.Context) {
	courseID, _ := strconv.Atoi(c.Param("id"))

	userID := 0
	role := ""
	if uid := middleware.GetUserID(c); uid != 0 {
		userID = uid
		if r, ok := c.Get("role"); ok {
			role = r.(string)
		}
	}

	input := getfull.Input{
		CourseID: courseID,
		UserID:   userID,
		Role:     role,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
