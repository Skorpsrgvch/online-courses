package admin

import (
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/admin/access"
	"github.com/gin-gonic/gin"
)

type GrantAccessHandler struct {
	usecase *access.Usecase
}

func NewGrantAccessHandler(usecase *access.Usecase) *GrantAccessHandler {
	return &GrantAccessHandler{usecase: usecase}
}

func (h *GrantAccessHandler) Handle(c *gin.Context) {
	var input access.Input

	// Привязка JSON к структуре
	if err := c.ShouldBindJSON(&input); err != nil {
		common.HandleError(c, common.HttpError("Неверный формат данных: user_id и course_id обязательны", http.StatusBadRequest))
		return
	}

	// Дополнительная валидация
	if input.UserID <= 0 || input.CourseID <= 0 {
		common.HandleError(c, common.HttpError("ID пользователя и курса должны быть положительными числами", http.StatusBadRequest))
		return
	}

	// Выполнение логики
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		// Маппинг ошибок на HTTP статусы
		status := http.StatusInternalServerError
		msg := err.Error()

		// Пример обработки конкретных ошибок
		if msg == "пользователь не найден" || msg == "курс не найден" {
			status = http.StatusNotFound
		}

		common.HandleError(c, common.HttpError(msg, status))
		return
	}

	// Если доступ уже есть (логическая ошибка, но не техническая)
	if !output.Success {
		c.JSON(http.StatusConflict, gin.H{
			"error": output.Message,
		})
		return
	}

	// Успех
	c.JSON(http.StatusOK, gin.H{
		"message":   output.Message,
		"user_id":   input.UserID,
		"course_id": input.CourseID,
	})
}
