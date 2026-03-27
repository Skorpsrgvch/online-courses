package review

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/review/approve"
	"github.com/gin-gonic/gin"
)

type ApproveHandler struct {
	usecase *approve.Usecase
}

func NewApproveHandler(usecase *approve.Usecase) *ApproveHandler {
	return &ApproveHandler{usecase: usecase}
}

func (h *ApproveHandler) Handle(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, common.HttpError("invalid review ID", http.StatusBadRequest))
		return
	}

	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	input := approve.Input{ReviewID: reviewID}
	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
