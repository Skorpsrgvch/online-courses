package list

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	CourseID int
	UserID   int
}

type Output struct {
	Reviews []ReviewDTO `json:"reviews"`
}

type ReviewDTO struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	CourseID        int       `json:"course_id"`
	Text            string    `json:"text"`
	Rating          int       `json:"rating"`
	Approved        bool      `json:"approved"`
	CreatedAt       time.Time `json:"created_at"`
	AuthorName      string    `json:"author_name"`
	RejectionReason string    `json:"rejection_reason,omitempty"`
}

type Usecase struct {
	reviewReader ReviewReader
}

type ReviewReader interface {
	GetApprovedReviewsByCourse(ctx context.Context, courseID int) ([]*domain.Review, error)
	GetByUserAndCourse(ctx context.Context, userID, courseID int) (*domain.Review, error)
}

func NewUsecase(reviewReader ReviewReader) (*Usecase, error) {
	if reviewReader == nil {
		return nil, errors.New("reviewReader is required")
	}
	return &Usecase{reviewReader: reviewReader}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Debug("ListReviews started", zap.Int("courseID", input.CourseID), zap.Int("userID", input.UserID))

	var reviews []*domain.Review

	if input.UserID > 0 {
		myReview, err := u.reviewReader.GetByUserAndCourse(ctx, input.UserID, input.CourseID)
		if err != nil {
			zap.L().Error("Failed to get user's review", zap.Int("userID", input.UserID), zap.Error(err))
			return nil, err
		}
		if myReview != nil {
			reviews = append(reviews, myReview)
			zap.L().Debug("Found user's own review", zap.Int("reviewID", myReview.ID))
		}
	} else {
		publicReviews, err := u.reviewReader.GetApprovedReviewsByCourse(ctx, input.CourseID)
		if err != nil {
			zap.L().Error("Failed to get public reviews", zap.Int("courseID", input.CourseID), zap.Error(err))
			return nil, err
		}
		reviews = publicReviews
		zap.L().Debug("Retrieved public reviews", zap.Int("count", len(reviews)))
	}

	var dtos []ReviewDTO
	for _, r := range reviews {
		dtos = append(dtos, ReviewDTO{
			ID:              r.ID,
			UserID:          r.UserID,
			CourseID:        r.CourseID,
			Text:            r.Text,
			Rating:          r.Rating,
			Approved:        r.Approved,
			CreatedAt:       r.CreatedAt,
			AuthorName:      r.AuthorName,
			RejectionReason: r.RejectionReason,
		})
	}

	zap.L().Info("Reviews listed successfully", zap.Int("courseID", input.CourseID), zap.Int("resultCount", len(dtos)))
	return &Output{Reviews: dtos}, nil
}
