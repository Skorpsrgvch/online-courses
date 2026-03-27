package domain

import (
	"fmt"
	"time"
)

type Review struct {
	ID        int
	UserID    int
	CourseID  int
	Text      string
	Rating    int // 1–5
	Approved  bool
	CreatedAt time.Time
}

func NewReview(text string, rating int, userID, courseID int) (*Review, error) {
	if text == "" {
		return nil, fmt.Errorf("review text is required")
	}
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}
	return &Review{
		Text:      text,
		Rating:    rating,
		UserID:    userID,
		CourseID:  courseID,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func RestoreReview(id, userID, courseID int, text string, rating int, approved bool, createdAt time.Time) *Review {
	return &Review{
		ID:        id,
		UserID:    userID,
		CourseID:  courseID,
		Text:      text,
		Rating:    rating,
		Approved:  approved,
		CreatedAt: createdAt,
	}
}
