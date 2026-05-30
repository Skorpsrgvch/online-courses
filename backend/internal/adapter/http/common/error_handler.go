package common

import (
	"errors"
	"net/http"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type httpError struct {
	msg    string
	status int
}

func (e httpError) Error() string { return e.msg }
func (e httpError) Status() int   { return e.status }

func HttpError(msg string, status int) error {
	return httpError{msg: msg, status: status}
}

func BadRequestError(msg string) error {
	return HttpError(msg, http.StatusBadRequest)
}

func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Проверка на наш кастомный HTTP error
	if httpErr, ok := err.(interface{ Status() int }); ok {
		if httpErr.Status() >= 500 {
			zap.L().Error("HTTP Error", zap.Int("status", httpErr.Status()), zap.String("message", err.Error()))
		}
		c.AbortWithStatusJSON(httpErr.Status(), gin.H{"error": err.Error()})
		return
	}

	// Сопоставление domain-ошибок
	switch {
	case errors.Is(err, domain.ErrUnauthorized):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})

	case errors.Is(err, domain.ErrTokenExpired):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})

	case errors.Is(err, domain.ErrTokenInvalid):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})

	case errors.Is(err, domain.ErrInvalidCredentials):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})

	case errors.Is(err, domain.ErrUserNotFound):
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
		zap.L().Error("Payment creation failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания платежа"})

	default:
		// Логируем неизвестную ошибку как ERROR
		zap.L().Error("Unhandled error in handler", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
	}
}
