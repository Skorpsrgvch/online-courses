package get_student_details

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	UserID int
}

type Output struct {
	Student *domain.StudentStat           `json:"student"`
	Courses []*domain.StudentCourseDetail `json:"courses"`
}

type StatsReader interface {
	GetStudentDetails(ctx context.Context, userID int) (*domain.StudentStat, []*domain.StudentCourseDetail, error)
}

type Usecase struct {
	statsReader StatsReader
}

func NewUsecase(statsReader StatsReader) (*Usecase, error) {
	if statsReader == nil {
		return nil, errors.New("statsReader is required")
	}
	return &Usecase{statsReader: statsReader}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Debug("GetStudentDetails started", zap.Int("userID", input.UserID))

	if input.UserID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	student, courses, err := u.statsReader.GetStudentDetails(ctx, input.UserID)
	if err != nil {
		zap.L().Error("Failed to get student details", zap.Int("userID", input.UserID), zap.Error(err))
		return nil, err
	}

	zap.L().Info("Student details fetched successfully", zap.Int("userID", input.UserID), zap.Int("coursesCount", len(courses)))
	return &Output{Student: student, Courses: courses}, nil
}
