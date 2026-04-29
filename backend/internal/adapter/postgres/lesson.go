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

// scanLesson — вспомогательная функция для сканирования строки с учетом private_key
func scanLesson(row *sql.Row) (*domain.Lesson, error) {
	var l domain.Lesson
	var pk sql.NullString

	err := row.Scan(
		&l.ID,
		&l.ModuleID,
		&l.Title,
		&l.Description,
		&l.VideoEmbedID,
		&pk,
		&l.Order,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrLessonNotFound
		}
		return nil, err
	}

	// Если ключ есть в БД, сохраняем его в структуру
	if pk.Valid {
		l.PrivateKey = &pk.String
	} else {
		l.PrivateKey = nil
	}

	return &l, nil
}

func (r *LessonRepo) GetByModuleID(ctx context.Context, moduleID int) ([]*domain.Lesson, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, module_id, title, description, video_embed_id, private_key, "order"
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
		var pk sql.NullString

		err := rows.Scan(
			&l.ID,
			&l.ModuleID,
			&l.Title,
			&l.Description,
			&l.VideoEmbedID,
			&pk,
			&l.Order,
		)
		if err != nil {
			return nil, err
		}

		if pk.Valid {
			l.PrivateKey = &pk.String
		} else {
			l.PrivateKey = nil
		}

		lessons = append(lessons, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lessons, nil
}

func (r *LessonRepo) GetByID(ctx context.Context, id int) (*domain.Lesson, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, module_id, title, description, video_embed_id, private_key, "order"
		 FROM lessons WHERE id = $1`,
		id,
	)

	return scanLesson(row)
}

func (r *LessonRepo) Save(ctx context.Context, lesson *domain.Lesson) error {
	query := `
		INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order")
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	// Подготовка значения private_key для записи (может быть NULL)
	var pkVal interface{}
	if lesson.PrivateKey != nil {
		pkVal = *lesson.PrivateKey
	} else {
		pkVal = nil
	}

	return r.db.QueryRowContext(ctx, query,
		lesson.ModuleID,
		lesson.Title,
		lesson.Description,
		lesson.VideoEmbedID,
		pkVal,
		lesson.Order,
	).Scan(&lesson.ID)
}

func (r *LessonRepo) Update(ctx context.Context, lesson *domain.Lesson) error {

	var pkVal interface{}
	if lesson.PrivateKey != nil {
		pkVal = *lesson.PrivateKey
	} else {
		pkVal = nil
	}

	_, err := r.db.ExecContext(ctx,
		`UPDATE lessons
		 SET title = $1, description = $2, video_embed_id = $3, private_key = $4, "order" = $5
		 WHERE id = $6`,
		lesson.Title,
		lesson.Description,
		lesson.VideoEmbedID,
		pkVal,
		lesson.Order,
		lesson.ID,
	)
	return err
}

func (r *LessonRepo) Delete(ctx context.Context, lessonID int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM lessons WHERE id = $1`,
		lessonID,
	)
	return err
}
