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
	query := `
		INSERT INTO user_progress (user_id, lesson_id, completed_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, lesson_id) 
		DO UPDATE SET completed_at = EXCLUDED.completed_at
	`
	_, err := r.db.ExecContext(ctx, query, userID, lessonID, time.Now().UTC())
	return err
}

// GetCourseProgress возвращает количество пройденных уроков для конкретного курса
func (r *ProgressRepo) GetCourseProgress(ctx context.Context, userID, courseID int) (completed int, total int, err error) {
	// Считаем всего уроков в курсе (через modules, так как в lessons нет course_id)
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM lessons l
		JOIN modules m ON l.module_id = m.id
		WHERE m.course_id = $1
	`, courseID).Scan(&total)
	if err != nil {
		return 0, 0, err
	}

	if total == 0 {
		return 0, 0, nil
	}

	// Считаем пройденные уроки пользователем в этом курсе
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM user_progress up
		JOIN lessons l ON up.lesson_id = l.id
		JOIN modules m ON l.module_id = m.id
		WHERE up.user_id = $1 AND m.course_id = $2
	`, userID, courseID).Scan(&completed)

	return completed, total, err
}

// GetCompletedLessonsCount возвращает количество пройденных уроков пользователем в курсе
// Исправлено: добавлен JOIN с modules для получения course_id
func (r *ProgressRepo) GetCompletedLessonsCount(ctx context.Context, userID, courseID int) (int, error) {
	query := `
		SELECT COUNT(DISTINCT up.lesson_id)
		FROM user_progress up
		JOIN lessons l ON l.id = up.lesson_id
		JOIN modules m ON l.module_id = m.id
		WHERE up.user_id = $1 AND m.course_id = $2
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID, courseID).Scan(&count)
	return count, err
}

// GetTotalLessonsCount возвращает общее количество уроков в курсе
// Исправлено: добавлен JOIN с modules для получения course_id
func (r *ProgressRepo) GetTotalLessonsCount(ctx context.Context, courseID int) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM lessons l
		JOIN modules m ON l.module_id = m.id
		WHERE m.course_id = $1
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, courseID).Scan(&count)
	return count, err
}

// IsLessonCompleted проверяет, прошел ли пользователь конкретный урок
func (r *ProgressRepo) IsLessonCompleted(ctx context.Context, userID, lessonID int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM user_progress WHERE user_id = $1 AND lesson_id = $2)",
		userID, lessonID,
	).Scan(&exists)
	return exists, err
}
