package postgres

import (
	"context"
	"database/sql"

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
		`SELECT id, module_id, title, description, video_embed_id, "order"
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
		err := rows.Scan(&l.ID, &l.ModuleID, &l.Title, &l.Description, &l.VideoEmbedID, &l.Order)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, &l)
	}
	return lessons, nil
}

func (r *LessonRepo) Save(ctx context.Context, lesson *domain.Lesson) error {
	query := `
		INSERT INTO lessons (module_id, title, description, video_embed_id, "order")
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		lesson.ModuleID, lesson.Title, lesson.Description, lesson.VideoEmbedID, lesson.Order,
	).Scan(&lesson.ID)
}
