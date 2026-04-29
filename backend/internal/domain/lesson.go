package domain

import "fmt"

type Lesson struct {
	ID           int     `json:"id"`
	ModuleID     int     `json:"module_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key,omitempty"`
	Order        int     `json:"order"`
}

func NewLesson(title, description, videoEmbedID string, moduleID, order int, privateKey *string) (*Lesson, error) {
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
		PrivateKey:   privateKey,
		Order:        order,
	}, nil
}

// RestoreLesson используется для восстановления объекта из базы данных
func RestoreLesson(id, moduleID, order int, title, description, videoEmbedID string, privateKey *string) *Lesson {
	return &Lesson{
		ID:           id,
		ModuleID:     moduleID,
		Title:        title,
		Description:  description,
		VideoEmbedID: videoEmbedID,
		PrivateKey:   privateKey,
		Order:        order,
	}
}
