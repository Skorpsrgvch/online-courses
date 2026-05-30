package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type AdminStatsRepo struct {
	db *sql.DB
}

func NewAdminStatsRepo(db *sql.DB) *AdminStatsRepo {
	return &AdminStatsRepo{db: db}
}

// GetTopStudents возвращает топ студентов по прогрессу
func (r *AdminStatsRepo) GetTopStudents(ctx context.Context, limit int) ([]*domain.StudentStat, error) {
	query := `
		SELECT 
			u.id,
			u.full_name,
			u.email,
			u.created_at as registered_at,
			COUNT(DISTINCT up.course_id) as total_courses,
			COUNT(DISTINCT CASE WHEN up.progress_percent = 100 THEN up.course_id END) as completed_courses,
			COALESCE(SUM(up.total_lessons), 0) as total_lessons,
			COALESCE(SUM(up.completed_lessons), 0) as completed_lessons,
			CASE 
				WHEN COALESCE(SUM(up.total_lessons), 0) = 0 THEN 0
				ELSE ROUND(COALESCE(SUM(up.completed_lessons), 0)::numeric / NULLIF(SUM(up.total_lessons), 0)::numeric * 100)
			END as progress_percent,
			MAX(up.last_activity_at) as last_activity_at
		FROM users u
		LEFT JOIN (
			SELECT 
				m.user_id,
				m.course_id,
				COUNT(DISTINCT l.id) as total_lessons,
				COUNT(DISTINCT p.lesson_id) as completed_lessons,
				CASE 
					WHEN COUNT(DISTINCT l.id) = COUNT(DISTINCT p.lesson_id) AND COUNT(DISTINCT l.id) > 0 THEN 100
					ELSE 0
				END as progress_percent,
				MAX(p.completed_at) as last_activity_at
			FROM user_purchases m
			JOIN modules mod ON mod.course_id = m.course_id
			JOIN lessons l ON l.module_id = mod.id
			LEFT JOIN user_progress p ON p.lesson_id = l.id AND p.user_id = m.user_id
			GROUP BY m.user_id, m.course_id
		) up ON u.id = up.user_id
		WHERE u.role = 'user'
		GROUP BY u.id, u.full_name, u.email, u.created_at
		ORDER BY progress_percent DESC, last_activity_at DESC NULLS LAST
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("query top students: %w", err)
	}
	defer rows.Close()

	var stats []*domain.StudentStat
	for rows.Next() {
		var s domain.StudentStat
		var lastAct sql.NullTime

		err := rows.Scan(
			&s.ID, &s.Name, &s.Email, &s.RegisteredAt,
			&s.TotalCourses, &s.CompletedCourses,
			&s.TotalLessons, &s.CompletedLessons,
			&s.ProgressPercent, &lastAct,
		)
		if err != nil {
			return nil, fmt.Errorf("scan student stat row: %w", err)
		}

		if lastAct.Valid {
			s.LastActivityAt = &lastAct.Time
		}

		stats = append(stats, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate student stat rows: %w", err)
	}

	zap.L().Debug("Top students fetched", zap.Int("count", len(stats)))
	return stats, nil
}

func (r *AdminStatsRepo) GetStudentDetails(ctx context.Context, userID int) (*domain.StudentStat, []*domain.StudentCourseDetail, error) {
	// 1. Получаем общую статистику пользователя
	var stat domain.StudentStat
	var lastAct sql.NullTime

	queryStat := `
		SELECT 
   		u.id, u.full_name, u.email, u.created_at,
    	COUNT(DISTINCT up.course_id),
    	COUNT(DISTINCT CASE WHEN cp.completed = cp.total THEN up.course_id END),
    	COALESCE(SUM(cp.total), 0),
    	COALESCE(SUM(cp.completed), 0),
    	CASE WHEN COALESCE(SUM(cp.total), 0) = 0 THEN 0 ELSE ROUND(COALESCE(SUM(cp.completed), 0)::numeric / NULLIF(SUM(cp.total), 0)::numeric * 100) END,
    	MAX(cp.last_activity)
		FROM users u
		LEFT JOIN user_purchases up ON u.id = up.user_id
		LEFT JOIN (
    	SELECT 
        m.user_id, m.course_id, 
        COUNT(DISTINCT l.id) as total,
        COUNT(DISTINCT prog.lesson_id) as completed,
        MAX(prog.completed_at) as last_activity
    	FROM user_purchases m
    	JOIN modules mod ON mod.course_id = m.course_id
    	JOIN lessons l ON l.module_id = mod.id
    	LEFT JOIN user_progress prog ON prog.lesson_id = l.id AND prog.user_id = m.user_id
    	GROUP BY m.user_id, m.course_id 
		) cp ON up.user_id = cp.user_id AND up.course_id = cp.course_id
		WHERE u.id = $1
		GROUP BY u.id
		`

	err := r.db.QueryRowContext(ctx, queryStat, userID).Scan(
		&stat.ID, &stat.Name, &stat.Email, &stat.RegisteredAt,
		&stat.TotalCourses, &stat.CompletedCourses,
		&stat.TotalLessons, &stat.CompletedLessons,
		&stat.ProgressPercent, &lastAct,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("user not found")
		}
		return nil, nil, fmt.Errorf("query student details: %w", err)
	}

	if lastAct.Valid {
		stat.LastActivityAt = &lastAct.Time
	}

	// 2. Получаем детали по каждому курсу
	queryCourses := `
		SELECT 
			c.id, c.title, up.purchased_at,
			COUNT(DISTINCT l.id) as total_lessons,
			COUNT(DISTINCT p.lesson_id) as completed_lessons,
			CASE WHEN COUNT(DISTINCT l.id) = 0 THEN 0 ELSE ROUND(COUNT(DISTINCT p.lesson_id)::numeric / COUNT(DISTINCT l.id)::numeric * 100) END,
			COALESCE(last_l.title, ''),
			last_p.completed_at
		FROM user_purchases up
		JOIN courses c ON c.id = up.course_id
		JOIN modules m ON m.course_id = c.id
		JOIN lessons l ON l.module_id = m.id
		LEFT JOIN user_progress p ON p.lesson_id = l.id AND p.user_id = up.user_id
		LEFT JOIN LATERAL (
			SELECT prog.lesson_id, prog.completed_at, les.title
			FROM user_progress prog
			JOIN lessons les ON les.id = prog.lesson_id
			JOIN modules mod2 ON mod2.id = les.module_id
			WHERE mod2.course_id = c.id AND prog.user_id = up.user_id
			ORDER BY prog.completed_at DESC
			LIMIT 1
		) last_p ON true
		LEFT JOIN lessons last_l ON last_l.id = last_p.lesson_id
		WHERE up.user_id = $1
		GROUP BY c.id, c.title, up.purchased_at, last_l.title, last_p.completed_at
		ORDER BY up.purchased_at DESC
	`

	rows, err := r.db.QueryContext(ctx, queryCourses, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("query student courses: %w", err)
	}
	defer rows.Close()

	var courses []*domain.StudentCourseDetail
	for rows.Next() {
		var c domain.StudentCourseDetail
		var lastActCourse sql.NullTime

		err := rows.Scan(
			&c.CourseID, &c.CourseTitle, &c.PurchasedAt,
			&c.TotalLessons, &c.CompletedLessons, &c.ProgressPercent,
			&c.LastLessonTitle, &lastActCourse,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("scan student course row: %w", err)
		}

		if lastActCourse.Valid {
			c.LastActivityAt = &lastActCourse.Time
		}

		courses = append(courses, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterate student course rows: %w", err)
	}

	zap.L().Debug("Student details fetched", zap.Int("userID", userID), zap.Int("coursesCount", len(courses)))
	return &stat, courses, nil
}
