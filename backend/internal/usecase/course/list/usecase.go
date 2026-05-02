package list

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Output struct {
	Courses []*domain.Course
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
	courses, err := u.courseLister.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return &Output{Courses: courses}, nil
}

// ИСПРАВЛЕНО: используем courseLister и возвращаем тот же тип, что и репозиторий
func (u *Usecase) ExecuteAdmin(ctx context.Context) ([]*domain.Course, error) {
	return u.courseLister.GetAllWithInactive(ctx)
}
