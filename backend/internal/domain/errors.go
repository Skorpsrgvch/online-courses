package domain

import "errors"

// Ошибки сущностей
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrCourseNotFound   = errors.New("course not found")
	ErrLessonNotFound   = errors.New("lesson not found")
	ErrModuleNotFound   = errors.New("module not found")
	ErrServiceNotFound  = errors.New("service not found")
	ErrReviewNotFound   = errors.New("review not found")
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrPurchaseNotFound = errors.New("purchase not found")
)

// Ошибки авторизации и прав доступа
var (
	ErrUnauthorized       = errors.New("unauthorized")
	ErrAccessDenied       = errors.New("access denied")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenInvalid       = errors.New("invalid token")
)

// Ошибки бизнес-логики
var (
	ErrReviewAlreadyExists   = errors.New("review already exists for this course")
	ErrCourseNotPurchased    = errors.New("course not purchased")
	ErrInvalidInput          = errors.New("invalid input data")
	ErrInvalidID             = errors.New("invalid ID format")
	ErrPaymentExpired        = errors.New("payment expired")
	ErrPaymentAlreadyPaid    = errors.New("payment already paid")
	ErrPaymentInvalidStatus  = errors.New("invalid payment status")
	ErrPaymentCreationFailed = errors.New("invalid create payment")
)
