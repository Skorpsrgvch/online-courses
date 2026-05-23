package admin

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	searchUC "github.com/Skorpsrgvch/online-courses/internal/usecase/admin/search"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SearchUsersHandler struct {
	usecase *searchUC.Usecase
}

func NewSearchUsersHandler(usecase *searchUC.Usecase) *SearchUsersHandler {
	return &SearchUsersHandler{usecase: usecase}
}

func (h *SearchUsersHandler) Handle(c *gin.Context) {
	var input searchUC.Input

	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Debug("Invalid JSON in search users", zap.Error(err))
		common.HandleError(c, common.HttpError("Неверный формат запроса", http.StatusBadRequest))
		return
	}

	if input.EmailQuery == "" {
		common.HandleError(c, common.HttpError("Поле email обязательно для поиска", http.StatusBadRequest))
		return
	}

	zap.L().Debug("Searching users", zap.String("query", input.EmailQuery))

	users, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		zap.L().Error("Search users failed", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	if len(users) == 0 {
		zap.L().Info("No users found for query", zap.String("query", input.EmailQuery))
	} else {
		zap.L().Info("Users found", zap.Int("count", len(users)))
	}

	c.JSON(http.StatusOK, users)
}
