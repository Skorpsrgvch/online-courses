package review

import (
	"net/http"
	"strconv"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/review/reject"
	"github.com/gin-gonic/gin"
)

type RejectHandler struct {
	usecase *reject.Usecase
}

func NewRejectHandler(usecase *reject.Usecase) *RejectHandler {
	return &RejectHandler{usecase: usecase}
}

func (h *RejectHandler) Handle(c *gin.Context) {
	if !middleware.RequireAdmin(c) {
		common.HandleError(c, domain.ErrAccessDenied)
		return
	}

	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.HandleError(c, common.HttpError("invalid review ID", http.StatusBadRequest))
		return
	}

	input := reject.Input{ReviewID: reviewID}
	if err := h.usecase.Execute(c.Request.Context(), input); err != nil {
		common.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
