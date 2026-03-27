package domain

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrCourseNotFound      = errors.New("course not found")
	ErrLessonNotFound      = errors.New("lesson not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrAccessDenied        = errors.New("access denied")
	ErrReviewAlreadyExists = errors.New("review already exists")
	ErrCourseNotPurchased  = errors.New("course not purchased")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrModuleNotFound      = errors.New("module not found")
)
