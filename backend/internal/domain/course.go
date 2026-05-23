package domain

import (
	"fmt"
	"time"
)

type Course struct {
	ID                int         `json:"id"`
	Title             string      `json:"title"`
	Description       string      `json:"description"`
	IsPublic          bool        `json:"is_public"`
	Price             int         `json:"price"`
	AuthorID          int         `json:"author_id"`
	IsActive          bool        `json:"is_active"`
	CoverImageURL     string      `json:"cover_image_url"`
	Contraindications string      `json:"contraindications"`
	Recommendations   string      `json:"recommendations"`
	TargetAudience    string      `json:"target_audience"`
	CourseBasis       string      `json:"course_basis"`
	ClassBasis        string      `json:"class_basis"`
	Bonuses           []BonusItem `json:"bonuses"`
	IsPurchased       bool        `json:"is_purchased,omitempty"`
	ProgressPercent   int         `json:"progress_percent"`
	CreatedAt         time.Time   `json:"created_at"`
}

type BonusItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func NewCourse(
	title, description string,
	isPublic bool,
	price int,
	authorID int,
	coverImageURL string,
	contraindications, recommendations, targetAudience, courseBasis, classBasis string,
) (*Course, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required: %w", ErrInvalidInput)
	}
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative: %w", ErrInvalidInput)
	}

	return &Course{
		Title:             title,
		Description:       description,
		IsPublic:          isPublic,
		Price:             price,
		AuthorID:          authorID,
		IsActive:          true,
		CoverImageURL:     coverImageURL,
		Contraindications: contraindications,
		Recommendations:   recommendations,
		TargetAudience:    targetAudience,
		CourseBasis:       courseBasis,
		ClassBasis:        classBasis,
		Bonuses:           make([]BonusItem, 0),
		CreatedAt:         time.Now().UTC(),
	}, nil
}

func RestoreCourse(
	id, authorID int, title, description string, isPublic bool, price int, isActive bool,
	coverImageURL string, createdAt time.Time,
	contraindications, recommendations, targetAudience, courseBasis, classBasis string,
	bonuses []BonusItem,
) *Course {
	return &Course{
		ID:                id,
		Title:             title,
		Description:       description,
		IsPublic:          isPublic,
		Price:             price,
		AuthorID:          authorID,
		IsActive:          isActive,
		CoverImageURL:     coverImageURL,
		Contraindications: contraindications,
		Recommendations:   recommendations,
		TargetAudience:    targetAudience,
		CourseBasis:       courseBasis,
		ClassBasis:        classBasis,
		Bonuses:           bonuses,
		CreatedAt:         createdAt,
	}
}
