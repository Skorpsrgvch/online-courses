package enroll_free

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type Input struct {
	UserID      int
	CourseID    int
	CoursePrice int
}

type Output struct {
	Message string `json:"message"`
}

type Purchaser interface {
	EnrollFree(ctx context.Context, userID, courseID, coursePrice int) error
}

type Usecase struct {
	purchaser Purchaser
}

func NewUsecase(purchaser Purchaser) (*Usecase, error) {
	if purchaser == nil {
		return nil, errors.New("purchaser dependency is required")
	}
	return &Usecase{
		purchaser: purchaser,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Debug("EnrollFree started", zap.Any("input", input))

	if input.CourseID <= 0 || input.UserID <= 0 {
		err := errors.New("invalid user or course ID")
		zap.L().Warn("EnrollFree validation failed", zap.Error(err))
		return nil, err
	}

	if err := u.purchaser.EnrollFree(ctx, input.UserID, input.CourseID, input.CoursePrice); err != nil {
		zap.L().Error("EnrollFree failed", zap.Error(err))
		return nil, err
	}

	zap.L().Info("User enrolled to free course",
		zap.Int("user_id", input.UserID),
		zap.Int("course_id", input.CourseID),
	)
	return &Output{Message: "Доступ к бесплатному курсу предоставлен"}, nil
}
