package access

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	UserID   int `json:"user_id"`
	CourseID int `json:"course_id"`
}

type Output struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Зависимости
type UserRepo interface {
	GetByID(ctx context.Context, id int) (*domain.User, error)
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
	// 1. Проверяем существование пользователя
	_, err := u.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// 2. Проверяем существование курса
	_, err = u.courseRepo.GetByID(ctx, input.CourseID)
	if err != nil {
		return nil, errors.New("курс не найден")
	}

	// 3. Выдаем доступ
	err = u.purchaseRepo.GrantAccess(ctx, input.UserID, input.CourseID)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentAlreadyPaid) {
			return &Output{
				Success: false,
				Message: "У пользователя уже есть доступ к этому курсу",
			}, nil
		}
		return nil, err
	}

	return &Output{
		Success: true,
		Message: "Доступ успешно предоставлен",
	}, nil
}
