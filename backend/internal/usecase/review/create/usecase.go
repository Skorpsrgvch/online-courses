package create

import (
	"context"
	"errors"
	"log"

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
	log.Printf("[INFO] CreateReviewUsecase.Execute: Start for UserID=%d, CourseID=%d", input.UserID, input.CourseID)

	// Проверяем существование пользователя
	log.Printf("[DEBUG] Checking user existence: ID=%d", input.UserID)
	if _, err := u.userFinder.GetUserByID(ctx, input.UserID); err != nil {
		log.Printf("[ERROR] User not found: %v", err)
		return err
	}

	// Проверяем существование курса
	log.Printf("[DEBUG] Checking course existence: ID=%d", input.CourseID)
	if _, err := u.courseFinder.GetByID(ctx, input.CourseID); err != nil {
		log.Printf("[ERROR] Course not found: %v", err)
		return err
	}

	// Создаём отзыв
	log.Printf("[DEBUG] Creating domain review object")
	review, err := domain.NewReview(input.Text, input.Rating, input.UserID, input.CourseID)
	if err != nil {
		log.Printf("[ERROR] Domain validation failed: %v", err)
		return err
	}

	log.Printf("[DEBUG] Saving review to repository")
	if err := u.reviewSaver.CreateReview(ctx, review); err != nil {
		log.Printf("[ERROR] Repository save failed: %v", err)
		return err
	}

	log.Printf("[INFO] CreateReviewUsecase.Execute: Success")
	return nil
}
