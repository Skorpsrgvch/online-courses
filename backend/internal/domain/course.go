package domain

import "fmt"

type Course struct {
	ID                int
	Title             string
	Description       string
	IsPublic          bool
	Price             int
	AuthorID          int
	IsActive          bool
	CoverImageURL     string
	Contraindications string      `json:"contraindications"` // Противопоказания
	Recommendations   string      `json:"recommendations"`
	TargetAudience    string      `json:"target_audience"`
	CourseBasis       string      `json:"course_basis"`
	ClassBasis        string      `json:"class_basis"`
	Bonuses           []BonusItem `json:"bonuses"` // Бонусы
	IsPurchased       bool        `json:"is_purchased,omitempty"`
}

type BonusItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"` // например: "gift", "book", "video"
}

// NewCourse создает новый курс с валидацией
func NewCourse(
	title, description string,
	isPublic bool,
	price int,
	authorID int,
	coverImageURL string,
	contraindications string,
	recommendations string,
	targetAudience string,
	courseBasis string,
	classBasis string,
) (*Course, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
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
		Bonuses:           []BonusItem{},
	}, nil
}

// RestoreCourse используется для восстановления объекта из базы данных
func RestoreCourse(
	id int, title, description string, isPublic bool, price int, authorID int, isActive bool,
	coverImageURL, contraindications, recommendations, targetAudience, courseBasis, classBasis string,
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
		Bonuses:           []BonusItem{},
	}
}
