package domain

import "fmt"

type Module struct {
	ID       int    `json:"id"`
	CourseID int    `json:"course_id"`
	Title    string `json:"title"`
	Order    int    `json:"order"`
}

func NewModule(title string, courseID, order int) (*Module, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required: %w", ErrInvalidInput)
	}
	return &Module{
		Title:    title,
		CourseID: courseID,
		Order:    order,
	}, nil
}

func RestoreModule(id, courseID, order int, title string) *Module {
	return &Module{
		ID:       id,
		CourseID: courseID,
		Title:    title,
		Order:    order,
	}
}
