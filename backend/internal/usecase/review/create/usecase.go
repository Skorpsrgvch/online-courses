package create

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	UserID   int
	CourseID int
	Text     string
	Rating   int
}

type Usecase struct {
	reviewSaver  ReviewSaver
	userFinder   UserFinder
	courseFinder CourseFinder
}

func NewUsecase(reviewSaver ReviewSaver, userFinder UserFinder, courseFinder CourseFinder) (*Usecase, error) {
	if reviewSaver == nil || userFinder == nil || courseFinder == nil {
		return nil, errors.New("all dependencies are required")
	}
	return &Usecase{
		reviewSaver:  reviewSaver,
		userFinder:   userFinder,
		courseFinder: courseFinder,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	// Проверяем существование пользователя
	if _, err := u.userFinder.GetByID(ctx, input.UserID); err != nil {
		return err
	}

	// Проверяем существование курса
	if _, err := u.courseFinder.GetByID(ctx, input.CourseID); err != nil {
		return err
	}

	// Создаём отзыв
	review, err := domain.NewReview(input.Text, input.Rating, input.UserID, input.CourseID)
	if err != nil {
		return err
	}

	return u.reviewSaver.CreateReview(ctx, review)
}
