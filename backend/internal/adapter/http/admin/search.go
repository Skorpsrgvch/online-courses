package admin

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	admin "github.com/Skorpsrgvch/online-courses/internal/usecase/admin/search"
	"github.com/gin-gonic/gin"
)

type SearchUsersHandler struct {
	usecase *admin.Usecase
}

func NewSearchUsersHandler(usecase *admin.Usecase) *SearchUsersHandler {
	return &SearchUsersHandler{usecase: usecase}
}

func (h *SearchUsersHandler) Handle(c *gin.Context) {
	var input admin.Input

	// Привязка JSON к структуре Input
	if err := c.ShouldBindJSON(&input); err != nil {
		common.HandleError(c, common.HttpError("Неверный формат запроса", http.StatusBadRequest))
		return
	}

	// Базовая валидация
	if input.EmailQuery == "" {
		common.HandleError(c, common.HttpError("Поле email обязательно для поиска", http.StatusBadRequest))
		return
	}

	// Выполнение логики
	users, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}
