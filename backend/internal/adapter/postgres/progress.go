package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type ProgressRepo struct {
	db *sql.DB
}

func NewProgressRepo(db *sql.DB) *ProgressRepo {
	return &ProgressRepo{db: db}
}

// GetCourseStartDate возвращает дату начала курса (дату покупки).
func (r *ProgressRepo) GetCourseStartDate(ctx context.Context, userID, courseID int) (*time.Time, error) {
	var purchasedAt time.Time
	err := r.db.QueryRowContext(ctx, `
		SELECT purchased_at FROM user_purchases 
		WHERE user_id = $1 AND course_id = $2
	`, userID, courseID).Scan(&purchasedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zap.L().Debug("Course purchase not found", zap.Int("user_id", userID), zap.Int("course_id", courseID))
			return nil, nil
		}
		zap.L().Error("Failed to get course start date", zap.Int("user_id", userID), zap.Int("course_id", courseID), zap.Error(err))
		return nil, fmt.Errorf("get course start date: %w", err)
	}

	zap.L().Debug("Course start date fetched", zap.Int("user_id", userID), zap.Int("course_id", courseID), zap.Time("started_at", purchasedAt))
	return &purchasedAt, nil
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
	if err != nil {
		zap.L().Error("Failed to mark lesson as completed", zap.Int("user_id", userID), zap.Int("lesson_id", lessonID), zap.Error(err))
		return fmt.Errorf("mark lesson completed: %w", err)
	}

	zap.L().Debug("Lesson marked as completed", zap.Int("user_id", userID), zap.Int("lesson_id", lessonID))
	return nil
}

// GetCourseProgress возвращает количество пройденных уроков для конкретного курса
func (r *ProgressRepo) GetCourseProgress(ctx context.Context, userID, courseID int) (completed int, total int, err error) {
	// Считаем всего уроков в курсе
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM lessons l
		JOIN modules m ON l.module_id = m.id
		WHERE m.course_id = $1
	`, courseID).Scan(&total)

	if err != nil {
		zap.L().Error("Failed to count total lessons", zap.Int("course_id", courseID), zap.Error(err))
		return 0, 0, fmt.Errorf("count total lessons: %w", err)
	}

	if total == 0 {
		zap.L().Debug("No lessons found in course", zap.Int("course_id", courseID))
		return 0, 0, nil
	}

	// Считаем пройденные уроки пользователем в этом курсе
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM user_progress up
		JOIN lessons l ON up.lesson_id = l.id
		JOIN modules m ON l.module_id = m.id
		WHERE up.user_id = $1 AND m.course_id = $2
	`, userID, courseID).Scan(&completed)

	if err != nil {
		zap.L().Error("Failed to count completed lessons", zap.Int("user_id", userID), zap.Int("course_id", courseID), zap.Error(err))
		return 0, 0, fmt.Errorf("count completed lessons: %w", err)
	}

	zap.L().Debug("Course progress fetched", zap.Int("user_id", userID), zap.Int("course_id", courseID), zap.Int("completed", completed), zap.Int("total", total))
	return completed, total, nil
}

// GetCompletedLessonsCount возвращает количество пройденных уроков пользователем в курсе
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
	if err != nil {
		zap.L().Error("Failed to count completed lessons", zap.Int("user_id", userID), zap.Int("course_id", courseID), zap.Error(err))
		return 0, fmt.Errorf("count completed lessons: %w", err)
	}

	zap.L().Debug("Completed lessons count fetched", zap.Int("user_id", userID), zap.Int("course_id", courseID), zap.Int("count", count))
	return count, nil
}

// GetTotalLessonsCount возвращает общее количество уроков в курсе
func (r *ProgressRepo) GetTotalLessonsCount(ctx context.Context, courseID int) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM lessons l
		JOIN modules m ON l.module_id = m.id
		WHERE m.course_id = $1
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, courseID).Scan(&count)
	if err != nil {
		zap.L().Error("Failed to count total lessons", zap.Int("course_id", courseID), zap.Error(err))
		return 0, fmt.Errorf("count total lessons: %w", err)
	}

	zap.L().Debug("Total lessons count fetched", zap.Int("course_id", courseID), zap.Int("count", count))
	return count, nil
}

// IsLessonCompleted проверяет, прошел ли пользователь конкретный урок
func (r *ProgressRepo) IsLessonCompleted(ctx context.Context, userID, lessonID int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM user_progress WHERE user_id = $1 AND lesson_id = $2)",
		userID, lessonID,
	).Scan(&exists)

	if err != nil {
		zap.L().Error("Failed to check lesson completion", zap.Int("user_id", userID), zap.Int("lesson_id", lessonID), zap.Error(err))
		return false, fmt.Errorf("check lesson completion: %w", err)
	}

	zap.L().Debug("Lesson completion status checked", zap.Int("user_id", userID), zap.Int("lesson_id", lessonID), zap.Bool("completed", exists))
	return exists, nil
}

func (r *PurchaseRepo) GetPurchaseDate(ctx context.Context, userID, courseID int) (*time.Time, error) {
	var purchasedAt time.Time
	err := r.db.QueryRowContext(ctx, `
		SELECT purchased_at FROM user_purchases 
		WHERE user_id = $1 AND course_id = $2
	`, userID, courseID).Scan(&purchasedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Покупка не найдена
		}
		zap.L().Error("Failed to get purchase date", zap.Int("userID", userID), zap.Int("courseID", courseID), zap.Error(err))
		return nil, err
	}

	zap.L().Debug("Purchase date fetched", zap.Int("userID", userID), zap.Int("courseID", courseID), zap.Time("date", purchasedAt))
	return &purchasedAt, nil
}
