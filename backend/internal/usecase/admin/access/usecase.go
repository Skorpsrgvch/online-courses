package access

import (
	"context"
	"errors"
	"fmt"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	UserID   int `json:"user_id"`
	CourseID int `json:"course_id"`
}

type Output struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type UserRepo interface {
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
}

type CourseRepo interface {
	GetByID(ctx context.Context, id int) (*domain.Course, error)
}

type PurchaseRepo interface {
	GrantAccess(ctx context.Context, userID, courseID int) error
}

type Usecase struct {
	userRepo     UserRepo
	courseRepo   CourseRepo
	purchaseRepo PurchaseRepo
}

func NewUsecase(userRepo UserRepo, courseRepo CourseRepo, purchaseRepo PurchaseRepo) (*Usecase, error) {
	if userRepo == nil || courseRepo == nil || purchaseRepo == nil {
		return nil, errors.New("all repositories are required")
	}

	return &Usecase{
		userRepo:     userRepo,
		courseRepo:   courseRepo,
		purchaseRepo: purchaseRepo,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Info("Granting access started", zap.Int("user_id", input.UserID), zap.Int("course_id", input.CourseID))

	// 1. Проверяем пользователя
	if _, err := u.userRepo.GetUserByID(ctx, input.UserID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			zap.L().Warn("User not found for access grant", zap.Int("user_id", input.UserID))
			return nil, errors.New("пользователь не найден")
		}
		zap.L().Error("Database error checking user", zap.Error(err))
		return nil, fmt.Errorf("check user failed: %w", err)
	}

	// 2. Проверяем курс
	if _, err := u.courseRepo.GetByID(ctx, input.CourseID); err != nil {
		if errors.Is(err, domain.ErrCourseNotFound) {
			zap.L().Warn("Course not found for access grant", zap.Int("course_id", input.CourseID))
			return nil, errors.New("курс не найден")
		}
		zap.L().Error("Database error checking course", zap.Error(err))
		return nil, fmt.Errorf("check course failed: %w", err)
	}

	// 3. Выдаем доступ
	if err := u.purchaseRepo.GrantAccess(ctx, input.UserID, input.CourseID); err != nil {
		zap.L().Error("Failed to grant access in repo", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Access granted successfully", zap.Int("user_id", input.UserID), zap.Int("course_id", input.CourseID))
	return &Output{
		Success: true,
		Message: "Доступ успешно предоставлен",
	}, nil
}
