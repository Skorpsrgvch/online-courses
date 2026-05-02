package update

import (
	"context"
	"errors"
	"log"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

func NewUsecase(repo CourseRepository) (*Usecase, error) {
	if repo == nil {
		return nil, errors.New("course repository is required")
	}
	return &Usecase{repo: repo}, nil
}

func (u *Usecase) Execute(ctx context.Context, course *domain.Course) error {
	log.Printf("[INFO] Usecase.Execute: starting update for course ID=%d", course.ID)

	if course == nil || course.ID == 0 {
		err := errors.New("invalid course data")
		log.Printf("[ERROR] Usecase.Execute: validation failed - %v", err)
		return err
	}

	if err := u.repo.Update(ctx, course); err != nil {
		log.Printf("[ERROR] Usecase.Execute: repo update failed - %v", err)
		return err
	}

	log.Printf("[INFO] Usecase.Execute: finished successfully for course ID=%d", course.ID)
	return nil
}

// UpdateStatus переключает статус активности курса
func (u *Usecase) UpdateStatus(ctx context.Context, id int, isActive bool) error {
	log.Printf("[INFO] Usecase.UpdateStatus: starting for course ID=%d, target IsActive=%v", id, isActive)

	if id == 0 {
		err := errors.New("invalid course ID")
		log.Printf("[ERROR] Usecase.UpdateStatus: validation failed - %v", err)
		return err
	}

	if err := u.repo.UpdateStatus(ctx, id, isActive); err != nil {
		log.Printf("[ERROR] Usecase.UpdateStatus: repo update failed - %v", err)
		return err
	}

	log.Printf("[INFO] Usecase.UpdateStatus: finished successfully for course ID=%d", id)
	return nil
}
