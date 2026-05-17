package domain

import (
	"fmt"
	"time"
)

type Review struct {
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

func NewReview(text string, rating int, userID, courseID int) (*Review, error) {
	if text == "" {
		return nil, fmt.Errorf("review text is required")
	}
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}
	return &Review{
		Text:       text,
		Rating:     rating,
		UserID:     userID,
		CourseID:   courseID,
		CreatedAt:  time.Now().UTC(),
		AuthorName: "",
	}, nil
}

func RestoreReview(id, userID, courseID int, text string, rating int, approved bool, createdAt time.Time, authorName string) *Review {
	return &Review{
		ID:         id,
		UserID:     userID,
		CourseID:   courseID,
		Text:       text,
		Rating:     rating,
		Approved:   approved,
		CreatedAt:  createdAt,
		AuthorName: authorName,
	}
}
