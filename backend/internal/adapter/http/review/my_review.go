package review

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	myReviewsUC "github.com/Skorpsrgvch/online-courses/internal/usecase/review/myreviews"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MyReviewsHandler struct {
	usecase *myReviewsUC.Usecase
}

func NewMyReviewsHandler(usecase *myReviewsUC.Usecase) *MyReviewsHandler {
	return &MyReviewsHandler{usecase: usecase}
}

func (h *MyReviewsHandler) Handle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, domain.ErrUnauthorized)
		return
	}

	zap.L().Debug("Fetching user's reviews", zap.Int("userID", userID))

	input := myReviewsUC.Input{UserID: userID}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Failed to fetch user's reviews", zap.Int("userID", userID), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
