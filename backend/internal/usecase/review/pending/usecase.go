package pending

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
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

type ReviewReader interface {
	GetPendingReviews(ctx context.Context) ([]*domain.Review, error)
}

type Usecase struct {
	reviewReader ReviewReader
}

func NewUsecase(reviewReader ReviewReader) (*Usecase, error) {
	if reviewReader == nil {
		return nil, errors.New("reviewReader is required")
	}
	return &Usecase{reviewReader: reviewReader}, nil
}

func (u *Usecase) Execute(ctx context.Context, _ Input) (*Output, error) {
	reviews, err := u.reviewReader.GetPendingReviews(ctx)
	if err != nil {
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

	return &Output{Reviews: dtos}, nil
}
