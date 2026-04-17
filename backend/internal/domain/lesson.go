package domain

import "fmt"

type LessonType string

const (
	LessonTypeVideo   LessonType = "video"
	LessonTypeArticle LessonType = "article"
)

type Lesson struct {
	ID             int
	ModuleID       int
	Title          string
	Description    string
	LessonType     LessonType // "video" | "article"
	VideoEmbedID   string     // Rutube ID (nullable для статей)
	ArticleContent string     // HTML/Markdown контент для статей
	Order          int
}

func NewLesson(title, description string, lessonType LessonType, videoEmbedID, articleContent string, moduleID, order int) (*Lesson, error) {
	if title == "" {
		return nil, fmt.Errorf("lesson title is required")
	}
	if lessonType == "" {
		lessonType = LessonTypeVideo
	}
	if lessonType != LessonTypeVideo && lessonType != LessonTypeArticle {
		return nil, fmt.Errorf("invalid lesson type: %s", lessonType)
	}
	if lessonType == LessonTypeVideo && videoEmbedID == "" {
		return nil, fmt.Errorf("video embed ID is required for video lessons")
	}
	if lessonType == LessonTypeArticle && articleContent == "" {
		return nil, fmt.Errorf("article content is required for article lessons")
	}
	return &Lesson{
		Title:          title,
		Description:    description,
		LessonType:     lessonType,
		VideoEmbedID:   videoEmbedID,
		ArticleContent: articleContent,
		ModuleID:       moduleID,
		Order:          order,
	}, nil
}

func RestoreLesson(id, moduleID, order int, title, description string, lessonType LessonType, videoEmbedID, articleContent string) *Lesson {
	return &Lesson{
		ID:             id,
		ModuleID:       moduleID,
		Title:          title,
		Description:    description,
		LessonType:     lessonType,
		VideoEmbedID:   videoEmbedID,
		ArticleContent: articleContent,
		Order:          order,
	}
}
