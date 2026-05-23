package myreviews

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	UserID int
}

type Output struct {
	Reviews []ReviewDTO `json:"reviews"`
}

type ReviewDTO struct {
	ID              int       `json:"id"`
	CourseID        int       `json:"course_id"`
	CourseTitle     string    `json:"course_title"`
	Text            string    `json:"text"`
	Rating          int       `json:"rating"`
	Approved        bool      `json:"approved"`
	RejectionReason string    `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type Usecase struct {
	reviewReader ReviewReader
}

type ReviewReader interface {
	GetByUserID(ctx context.Context, userID int) ([]*domain.Review, error)
}

func NewUsecase(reviewReader ReviewReader) (*Usecase, error) {
	if reviewReader == nil {
		return nil, errors.New("reviewReader is required")
	}
	return &Usecase{reviewReader: reviewReader}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Debug("GetMyReviews started", zap.Int("userID", input.UserID))

	if input.UserID <= 0 {
		err := errors.New("invalid user ID")
		zap.L().Warn("Validation failed", zap.Error(err))
		return nil, err
	}

	reviews, err := u.reviewReader.GetByUserID(ctx, input.UserID)
	if err != nil {
		zap.L().Error("Failed to get user reviews", zap.Int("userID", input.UserID), zap.Error(err))
		return nil, err
	}

	var dtos []ReviewDTO
	for _, r := range reviews {
		dtos = append(dtos, ReviewDTO{
			ID:              r.ID,
			CourseID:        r.CourseID,
			CourseTitle:     r.CourseTitle,
			Text:            r.Text,
			Rating:          r.Rating,
			Approved:        r.Approved,
			RejectionReason: r.RejectionReason,
			CreatedAt:       r.CreatedAt,
		})
	}

	zap.L().Info("My reviews retrieved", zap.Int("userID", input.UserID), zap.Int("count", len(dtos)))
	return &Output{Reviews: dtos}, nil
}
