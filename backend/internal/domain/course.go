package domain

import "fmt"

type Course struct {
	ID            int
	Title         string
	Description   string
	IsPublic      bool // true = бесплатный, false = платный
	Price         int  // стоимость в рублях (может быть 0)
	AuthorID      int
	IsActive      bool
	CoverImageURL string // URL обложки курса (nullable)
}

func NewCourse(title, description string, isPublic bool, price int, authorID int, coverImageURL string) (*Course, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}
	return &Course{
		Title:         title,
		Description:   description,
		IsPublic:      isPublic,
		Price:         price,
		AuthorID:      authorID,
		IsActive:      true,
		CoverImageURL: coverImageURL,
	}, nil
}

func RestoreCourse(id int, title, description string, isPublic bool, price int, authorID int, isActive bool, coverImageURL string) *Course {
	return &Course{
		ID:            id,
		Title:         title,
		Description:   description,
		IsPublic:      isPublic,
		Price:         price,
		AuthorID:      authorID,
		IsActive:      isActive,
		CoverImageURL: coverImageURL,
	}
}
