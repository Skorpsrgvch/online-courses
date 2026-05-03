package common

import (
	"errors"
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
	if httpErr, ok := err.(interface{ Status() int }); ok {
		c.AbortWithStatusJSON(httpErr.Status(), gin.H{"error": err.Error()})
		return
	}

	// Сопоставление domain-ошибок
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	case errors.Is(err, domain.ErrAccessDenied):
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
	case errors.Is(err, domain.ErrCourseNotPurchased):
		c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{"error": "course not purchased"})
	case errors.Is(err, domain.ErrPaymentNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "payment not found"})
	case errors.Is(err, domain.ErrPaymentAlreadyPaid):
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "course already purchased"})
	case errors.Is(err, domain.ErrPaymentExpired):
		c.AbortWithStatusJSON(http.StatusGone, gin.H{"error": "payment expired"})
	case errors.Is(err, domain.ErrPaymentInvalidStatus):
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid payment status"})
	case errors.Is(err, domain.ErrPaymentCreationFailed):
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment"})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
