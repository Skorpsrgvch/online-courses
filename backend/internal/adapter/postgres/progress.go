package postgres

import (
	"context"
	"database/sql"
	"time"
)

type ProgressRepo struct {
	db *sql.DB
}

func NewProgressRepo(db *sql.DB) *ProgressRepo {
	return &ProgressRepo{db: db}
}

// MarkCompleted отмечает урок как пройденный
func (r *ProgressRepo) MarkCompleted(ctx context.Context, userID, lessonID int) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_progress (user_id, lesson_id, completed_at)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, lesson_id) DO NOTHING`,
		userID, lessonID, time.Now().UTC(),
	)
	return err
}

// IsCompleted проверяет, пройден ли урок
func (r *ProgressRepo) IsCompleted(ctx context.Context, userID, lessonID int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM user_progress WHERE user_id = $1 AND lesson_id = $2)`,
		userID, lessonID,
	).Scan(&exists)
	return exists, err
}

// GetCompletedLessonsCount возвращает количество пройденных уроков в модулях указанного курса
func (r *ProgressRepo) GetCompletedLessonsCount(ctx context.Context, userID, courseID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT up.lesson_id)
		 FROM user_progress up
		 JOIN lessons l ON up.lesson_id = l.id
		 JOIN modules m ON l.module_id = m.id
		 WHERE up.user_id = $1 AND m.course_id = $2`,
		userID, courseID,
	).Scan(&count)
	return count, err
}

// GetTotalLessonsCount возвращает общее количество уроков в курсе
func (r *ProgressRepo) GetTotalLessonsCount(ctx context.Context, courseID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM lessons l
		 JOIN modules m ON l.module_id = m.id
		 WHERE m.course_id = $1`,
		courseID,
	).Scan(&count)
	return count, err
}
