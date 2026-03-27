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

func HttpError(msg string, status int) error {
	return httpError{msg: msg, status: status}
}

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
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
