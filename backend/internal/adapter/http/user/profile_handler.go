package user

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	userprofile "github.com/Skorpsrgvch/online-courses/internal/usecase/user/profile"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	usecase *userprofile.Usecase
}

func NewProfileHandler(usecase *userprofile.Usecase) *ProfileHandler {
	return &ProfileHandler{usecase: usecase}
}

func (h *ProfileHandler) Handle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, domain.ErrUnauthorized)
		return
	}

	input := userprofile.Input{UserID: userID}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
