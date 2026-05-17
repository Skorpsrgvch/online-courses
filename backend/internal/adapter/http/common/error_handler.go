package common

import (
	"errors"
	"log"
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/gin-gonic/gin"
)

type httpError struct {
	msg    string
	status int
}

func (e httpError) Error() string { return e.msg }
func (e httpError) Status() int   { return e.status }

// HttpError создает новую HTTP ошибку
func HttpError(msg string, status int) error {
	return httpError{msg: msg, status: status}
}

// BadRequestError удобная обертка для ошибок 400
func BadRequestError(msg string) error {
	return HttpError(msg, http.StatusBadRequest)
}

// HandleError централизованная обработка ошибок
func HandleError(c *gin.Context, err error) {
	// 1. Проверяем, является ли ошибка уже готовой HTTP ошибкой (наш кастомный тип)
	if httpErr, ok := err.(interface{ Status() int }); ok {
		c.AbortWithStatusJSON(httpErr.Status(), gin.H{"error": err.Error()})
		return
	}

	// 2. Сопоставление domain-ошибок
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		// Это самый важный кейс для логина
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})

	case errors.Is(err, domain.ErrUserNotFound):
		// На случай если где-то отдельно пробрасывается эта ошибка
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})

	case errors.Is(err, domain.ErrAccessDenied):
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Доступ закрыт"})

	case errors.Is(err, domain.ErrCourseNotPurchased):
		c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{"error": "Курс не оплачен"})

	case errors.Is(err, domain.ErrPaymentNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Платеж не найден"})

	case errors.Is(err, domain.ErrPaymentAlreadyPaid):
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Курс уже приобретен"})

	case errors.Is(err, domain.ErrPaymentExpired):
		c.AbortWithStatusJSON(http.StatusGone, gin.H{"error": "Срок платежа истек"})

	case errors.Is(err, domain.ErrPaymentInvalidStatus):
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Некорректный статус платежа"})

	case errors.Is(err, domain.ErrPaymentCreationFailed):
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания платежа"})

	default:
		// Логируем неизвестную ошибку для разработчика
		log.Printf("Unhandled error in handler: %v", err)
		// Клиенту показываем общую фразу, чтобы не светить детали реализации
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
	}
}
