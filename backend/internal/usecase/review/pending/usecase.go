package pending

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct{}

type ReviewDTO struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	CourseID    int       `json:"course_id"`
	Text        string    `json:"text"`
	Rating      int       `json:"rating"`
	Approved    bool      `json:"approved"`
	CreatedAt   time.Time `json:"created_at"`
	AuthorName  string    `json:"author_name"`
	CourseTitle string    `json:"course_title"`
}

type Output struct {
	Reviews []ReviewDTO `json:"reviews"`
}

type Usecase struct {
	reviewReader ReviewReader
}

type ReviewReader interface {
	GetPendingReviews(ctx context.Context) ([]*domain.Review, error)
}

func NewUsecase(reviewReader ReviewReader) (*Usecase, error) {
	if reviewReader == nil {
		return nil, errors.New("reviewReader is required")
	}
	return &Usecase{reviewReader: reviewReader}, nil
}

func (u *Usecase) Execute(ctx context.Context, _ Input) (*Output, error) {
	zap.L().Debug("GetPendingReviews started")

	reviews, err := u.reviewReader.GetPendingReviews(ctx)
	if err != nil {
		zap.L().Error("Failed to get pending reviews", zap.Error(err))
		return nil, err
	}

	var dtos []ReviewDTO
	for _, r := range reviews {
		dtos = append(dtos, ReviewDTO{
			ID:          r.ID,
			UserID:      r.UserID,
			CourseID:    r.CourseID,
			Text:        r.Text,
			Rating:      r.Rating,
			Approved:    r.Approved,
			CreatedAt:   r.CreatedAt,
			AuthorName:  r.AuthorName,
			CourseTitle: r.CourseTitle,
		})
	}

	zap.L().Info("Pending reviews retrieved", zap.Int("count", len(dtos)))
	return &Output{Reviews: dtos}, nil
}
