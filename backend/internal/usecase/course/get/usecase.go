package get

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	CourseID int
	UserID   int // 0 — если неавторизован
}

type Output struct {
	Course *domain.Course
}

type Usecase struct {
	courseReader    CourseReader
	purchaseChecker PurchaseChecker
}

func NewUsecase(courseReader CourseReader, purchaseChecker PurchaseChecker) (*Usecase, error) {
	if courseReader == nil || purchaseChecker == nil {
		return nil, errors.New("dependencies required")
	}
	return &Usecase{
		courseReader:    courseReader,
		purchaseChecker: purchaseChecker,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	course, err := u.courseReader.GetByID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	// Бесплатный курс — доступен всем
	if course.IsPublic {
		return &Output{Course: course}, nil
	}

	// Платный курс — только авторизованным
	if input.UserID == 0 {
		return nil, domain.ErrAccessDenied
	}

	// Проверяем покупку
	purchased, err := u.purchaseChecker.HasPurchased(ctx, input.UserID, course.ID)
	if err != nil {
		return nil, err
	}
	if !purchased {
		return nil, domain.ErrCourseNotPurchased
	}

	return &Output{Course: course}, nil
}
