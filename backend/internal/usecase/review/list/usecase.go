package list

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	CourseID int
}

type Output struct {
	Reviews []ReviewDTO `json:"reviews"`
}

type ReviewDTO struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	CourseID   int       `json:"course_id"`
	Text       string    `json:"text"`
	Rating     int       `json:"rating"`
	Approved   bool      `json:"approved"`
	CreatedAt  time.Time `json:"created_at"`
	AuthorName string    `json:"author_name"` // 1. Добавили поле
}

type ReviewReader interface {
	GetApprovedReviewsByCourse(ctx context.Context, courseID int) ([]*domain.Review, error)
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

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	reviews, err := u.reviewReader.GetApprovedReviewsByCourse(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	var dtos []ReviewDTO
	for _, r := range reviews {
		dtos = append(dtos, ReviewDTO{
			ID:         r.ID,
			UserID:     r.UserID,
			CourseID:   r.CourseID,
			Text:       r.Text,
			Rating:     r.Rating,
			Approved:   r.Approved,
			CreatedAt:  r.CreatedAt,
			AuthorName: r.AuthorName, // 2. Копируем имя автора из доменной модели
		})
	}

	return &Output{Reviews: dtos}, nil
}
