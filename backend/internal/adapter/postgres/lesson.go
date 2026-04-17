// internal/adapter/postgres/lesson.go
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type LessonRepo struct {
	db *sql.DB
}

func NewLessonRepo(db *sql.DB) *LessonRepo {
	return &LessonRepo{db: db}
}

func (r *LessonRepo) GetByModuleID(ctx context.Context, moduleID int) ([]*domain.Lesson, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, module_id, title, description, lesson_type, 
		        COALESCE(video_embed_id, ''), 
		        COALESCE(article_content, ''), 
		        "order"
		 FROM lessons WHERE module_id = $1 ORDER BY "order"`,
		moduleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []*domain.Lesson
	for rows.Next() {
		var l domain.Lesson
		var lessonType string
		err := rows.Scan(&l.ID, &l.ModuleID, &l.Title, &l.Description, &lessonType, &l.VideoEmbedID, &l.ArticleContent, &l.Order)
		if err != nil {
			return nil, err
		}
		l.LessonType = domain.LessonType(lessonType)
		lessons = append(lessons, &l)
	}
	return lessons, nil
}

func (r *LessonRepo) GetByID(ctx context.Context, id int) (*domain.Lesson, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, module_id, title, description, lesson_type, 
		        COALESCE(video_embed_id, ''), 
		        COALESCE(article_content, ''), 
		        "order"
		 FROM lessons WHERE id = $1`,
		id,
	)

	var l domain.Lesson
	var lessonType string
	err := row.Scan(&l.ID, &l.ModuleID, &l.Title, &l.Description, &lessonType, &l.VideoEmbedID, &l.ArticleContent, &l.Order)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrLessonNotFound
		}
		return nil, err
	}
	l.LessonType = domain.LessonType(lessonType)
	return &l, nil
}

func (r *LessonRepo) Save(ctx context.Context, lesson *domain.Lesson) error {
	query := `
		INSERT INTO lessons (module_id, title, description, lesson_type, video_embed_id, article_content, "order")
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''), $7)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		lesson.ModuleID,
		lesson.Title,
		lesson.Description,
		string(lesson.LessonType),
		lesson.VideoEmbedID,
		lesson.ArticleContent,
		lesson.Order,
	).Scan(&lesson.ID)
}

// Update обновляет урок
func (r *LessonRepo) Update(ctx context.Context, lesson *domain.Lesson) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE lessons
		 SET title = $1, description = $2, lesson_type = $3, 
		     video_embed_id = NULLIF($4, ''), article_content = $5, "order" = $6
		 WHERE id = $7`,
		lesson.Title,
		lesson.Description,
		string(lesson.LessonType),
		lesson.VideoEmbedID,
		lesson.ArticleContent,
		lesson.Order,
		lesson.ID,
	)
	return err
}

// Delete удаляет урок
func (r *LessonRepo) Delete(ctx context.Context, lessonID int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM lessons WHERE id = $1`,
		lessonID,
	)
	return err
}
