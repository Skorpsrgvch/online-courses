package domain

import "fmt"

type Lesson struct {
	ID           int
	ModuleID     int
	Title        string
	Description  string
	VideoEmbedID string // Rutube ID
	Order        int
}

func NewLesson(title, description, videoEmbedID string, moduleID, order int) (*Lesson, error) {
	if title == "" {
		return nil, fmt.Errorf("lesson title is required")
	}
	if videoEmbedID == "" {
		return nil, fmt.Errorf("video embed ID is required")
	}
	return &Lesson{
		Title:        title,
		Description:  description,
		VideoEmbedID: videoEmbedID,
		ModuleID:     moduleID,
		Order:        order,
	}, nil
}

func RestoreLesson(id, moduleID, order int, title, description, videoEmbedID string) *Lesson {
	return &Lesson{
		ID:           id,
		ModuleID:     moduleID,
		Title:        title,
		Description:  description,
		VideoEmbedID: videoEmbedID,
		Order:        order,
	}
}
