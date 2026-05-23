package update

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type CourseRepository interface {
	Update(ctx context.Context, course *domain.Course) error
	UpdateStatus(ctx context.Context, id int, isActive bool) error
}

type Usecase struct {
	repo CourseRepository
}

func NewUsecase(repo CourseRepository) (*Usecase, error) {
	if repo == nil {
		return nil, errors.New("course repository is required")
	}
	return &Usecase{
		repo: repo,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, course *domain.Course) error {
	zap.L().Info("Updating course", zap.Int("course_id", course.ID))

	if course == nil || course.ID == 0 {
		err := errors.New("invalid course data")
		zap.L().Warn("Update validation failed", zap.Error(err))
		return err
	}

	if err := u.repo.Update(ctx, course); err != nil {
		zap.L().Error("Course update failed", zap.Int("course_id", course.ID), zap.Error(err))
		return err
	}

	zap.L().Info("Course updated successfully", zap.Int("course_id", course.ID))
	return nil
}

func (u *Usecase) UpdateStatus(ctx context.Context, id int, isActive bool) error {
	zap.L().Info("Updating course status", zap.Int("course_id", id), zap.Bool("is_active", isActive))

	if id == 0 {
		err := errors.New("invalid course ID")
		zap.L().Warn("Status update validation failed", zap.Error(err))
		return err
	}

	if err := u.repo.UpdateStatus(ctx, id, isActive); err != nil {
		zap.L().Error("Course status update failed", zap.Int("course_id", id), zap.Error(err))
		return err
	}

	zap.L().Info("Course status updated", zap.Int("course_id", id))
	return nil
}
