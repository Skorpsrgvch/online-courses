package create

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
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

type ReviewSaver interface {
	CreateReview(ctx context.Context, review *domain.Review) error
}

type UserFinder interface {
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
}

type CourseFinder interface {
	GetByID(ctx context.Context, id int) (*domain.Course, error)
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
	zap.L().Debug("CreateReview started", zap.Int("userID", input.UserID), zap.Int("courseID", input.CourseID))

	if _, err := u.userFinder.GetUserByID(ctx, input.UserID); err != nil {
		zap.L().Error("User not found", zap.Int("userID", input.UserID), zap.Error(err))
		return err
	}

	if _, err := u.courseFinder.GetByID(ctx, input.CourseID); err != nil {
		zap.L().Error("Course not found", zap.Int("courseID", input.CourseID), zap.Error(err))
		return err
	}

	review, err := domain.NewReview(input.Text, input.Rating, input.UserID, input.CourseID)
	if err != nil {
		zap.L().Warn("Domain validation failed", zap.Error(err))
		return err
	}

	if err := u.reviewSaver.CreateReview(ctx, review); err != nil {
		zap.L().Error("Failed to save review", zap.Error(err))
		return err
	}

	zap.L().Info("Review created successfully", zap.Int("reviewID", review.ID), zap.Int("courseID", input.CourseID))
	return nil
}
