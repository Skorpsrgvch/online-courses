package payment

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/payment/list"
	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	usecase *list.UseCase
}

func NewListHandler(usecase *list.UseCase) *ListHandler {
	return &ListHandler{usecase: usecase}
}

func (h *ListHandler) Handle(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		common.HandleError(c, common.BadRequestError("user_id is required"))
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		common.HandleError(c, common.BadRequestError("invalid user_id format"))
		return
	}

	input := list.Input{
		UserID: userID,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
