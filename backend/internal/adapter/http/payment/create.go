package payment

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/payment/create"
	"github.com/gin-gonic/gin"
)

type CreateHandler struct {
	usecase *create.UseCase
}

func NewCreateHandler(usecase *create.UseCase) *CreateHandler {
	return &CreateHandler{usecase: usecase}
}

func (h *CreateHandler) Handle(c *gin.Context) {
	var input create.Input
	// Парсим только return_url и course_id из тела
	if err := c.ShouldBindJSON(&input); err != nil {
		common.HandleError(c, common.HttpError("invalid request body", http.StatusBadRequest))
		return
	}

	// Получаем UserID из авторизованного контекста
	userID := middleware.GetUserID(c)
	if userID == 0 {
		common.HandleError(c, common.HttpError("unauthorized", http.StatusUnauthorized))
		return
	}
	input.UserID = userID // Принудительно ставим ID из токена

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		// Обработка ошибок (уже купил, нет курса и т.д.)
		common.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, output)
}
