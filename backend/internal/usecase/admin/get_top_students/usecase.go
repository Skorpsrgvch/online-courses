package get_top_students

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	Limit int
}

type Output struct {
	Students []*domain.StudentStat `json:"students"`
}

type StatsReader interface {
	GetTopStudents(ctx context.Context, limit int) ([]*domain.StudentStat, error)
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
	zap.L().Debug("GetTopStudents started", zap.Int("limit", input.Limit))

	limit := input.Limit
	if limit <= 0 {
		limit = 20 // Дефолтное значение
	}

	students, err := u.statsReader.GetTopStudents(ctx, limit)
	if err != nil {
		zap.L().Error("Failed to get top students", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Top students fetched successfully", zap.Int("count", len(students)))
	return &Output{Students: students}, nil
}
