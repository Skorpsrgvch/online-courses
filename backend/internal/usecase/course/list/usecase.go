package list

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Output struct {
	Courses []*domain.Course
}

type CourseLister interface {
	ListAll(ctx context.Context) ([]*domain.Course, error)
	GetAllWithInactive(ctx context.Context) ([]*domain.Course, error)
}

type Usecase struct {
	courseLister CourseLister
}

func NewUsecase(courseLister CourseLister) (*Usecase, error) {
	if courseLister == nil {
		return nil, errors.New("courseLister is required")
	}
	return &Usecase{
		courseLister: courseLister,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context) (*Output, error) {
	zap.L().Debug("Listing active courses")
	courses, err := u.courseLister.ListAll(ctx)
	if err != nil {
		zap.L().Error("Failed to list courses", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Courses listed successfully", zap.Int("count", len(courses)))
	return &Output{Courses: courses}, nil
}

func (u *Usecase) ExecuteAdmin(ctx context.Context) ([]*domain.Course, error) {
	zap.L().Debug("Listing all courses (admin)")
	courses, err := u.courseLister.GetAllWithInactive(ctx)
	if err != nil {
		zap.L().Error("Failed to list admin courses", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Admin courses listed", zap.Int("count", len(courses)))
	return courses, nil
}
