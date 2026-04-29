package course

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
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
	idStr := c.Param("id")
	courseID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID курса"})
		return
	}

	userID := 0
	role := ""

	if uid, exists := c.Get("user_id"); exists && uid != nil {
		if idVal, ok := uid.(int); ok {
			userID = idVal
		}
	}

	if r, exists := c.Get("role"); exists && r != nil {
		if roleVal, ok := r.(string); ok {
			role = roleVal
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
