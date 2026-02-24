package domain

import "fmt"

type Module struct {
	ID       int
	CourseID int
	Title    string
	Order    int
}

func NewModule(title string, courseID, order int) (*Module, error) {
	if title == "" {
		return nil, fmt.Errorf("module title is required")
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
