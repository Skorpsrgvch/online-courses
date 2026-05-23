package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type CourseRepo struct {
	db *sql.DB
}

func NewCourseRepo(db *sql.DB) *CourseRepo {
	return &CourseRepo{db: db}
}

// scanCourse универсальная функция для сканирования строки курса из Row или Rows
func scanCourse(scanner interface{}) (*domain.Course, error) {
	var c domain.Course
	var bonusesRaw []byte
	var coverURL, contraindications, recommendations, targetAudience, courseBasis, classBasis sql.NullString

	var err error
	switch v := scanner.(type) {
	case *sql.Row:
		err = v.Scan(
			&c.ID, &c.Title, &c.Description, &c.IsPublic, &c.Price, &c.AuthorID, &c.IsActive,
			&coverURL, &contraindications, &recommendations, &targetAudience, &courseBasis, &classBasis,
			&bonusesRaw,
		)
	case *sql.Rows:
		err = v.Scan(
			&c.ID, &c.Title, &c.Description, &c.IsPublic, &c.Price, &c.AuthorID, &c.IsActive,
			&coverURL, &contraindications, &recommendations, &targetAudience, &courseBasis, &classBasis,
			&bonusesRaw,
		)
	default:
		return nil, errors.New("invalid scanner type")
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCourseNotFound
		}
		return nil, err
	}

	// Обработка NULL полей
	if coverURL.Valid {
		c.CoverImageURL = coverURL.String
	}
	if contraindications.Valid {
		c.Contraindications = contraindications.String
	}
	if recommendations.Valid {
		c.Recommendations = recommendations.String
	}
	if targetAudience.Valid {
		c.TargetAudience = targetAudience.String
	}
	if courseBasis.Valid {
		c.CourseBasis = courseBasis.String
	}
	if classBasis.Valid {
		c.ClassBasis = classBasis.String
	}

	// Парсинг JSON бонусов
	c.Bonuses = make([]domain.BonusItem, 0)
	if len(bonusesRaw) > 0 {
		if err := json.Unmarshal(bonusesRaw, &c.Bonuses); err != nil {
			zap.L().Warn("Failed to unmarshal course bonuses",
				zap.Int("course_id", c.ID),
				zap.Error(err),
			)
			// Возвращаем пустой слайс вместо ошибки, чтобы не ломать поток
			c.Bonuses = make([]domain.BonusItem, 0)
		}
	}

	return &c, nil
}

func (r *CourseRepo) GetByID(ctx context.Context, id int) (*domain.Course, error) {
	query := `
		SELECT id, title, description, is_public, price, author_id, is_active,
		       COALESCE(cover_image_url, ''), 
		       COALESCE(contraindications, ''), 
		       COALESCE(recommendations, ''),
		       COALESCE(target_audience, ''),
		       COALESCE(course_basis, ''),
		       COALESCE(class_basis, ''),
		       COALESCE(bonuses, '[]'::jsonb)
		FROM courses
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return scanCourse(row)
}

func (r *CourseRepo) ListAll(ctx context.Context) ([]*domain.Course, error) {
	query := `
		SELECT id, title, description, is_public, price, author_id, is_active,
		       COALESCE(cover_image_url, ''),
		       COALESCE(contraindications, ''),
		       COALESCE(recommendations, ''),
		       COALESCE(target_audience, ''),
		       COALESCE(course_basis, ''),
		       COALESCE(class_basis, ''),
		       COALESCE(bonuses, '[]'::jsonb)
		FROM courses
		WHERE is_active = true
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		zap.L().Error("Failed to execute ListAll query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	courses := make([]*domain.Course, 0)
	for rows.Next() {
		course, err := scanCourse(rows)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *CourseRepo) GetAllWithInactive(ctx context.Context) ([]*domain.Course, error) {
	query := `
		SELECT id, title, description, is_public, price, author_id, is_active,
		       COALESCE(cover_image_url, ''), 
		       COALESCE(contraindications, ''), 
		       COALESCE(recommendations, ''),
		       COALESCE(target_audience, ''),
		       COALESCE(course_basis, ''),
		       COALESCE(class_basis, ''),
		       COALESCE(bonuses, '[]'::jsonb)
		FROM courses
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := make([]*domain.Course, 0)
	for rows.Next() {
		course, err := scanCourse(rows)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, rows.Err()
}

func (r *CourseRepo) Save(ctx context.Context, course *domain.Course) error {
	bonusesJSON, err := json.Marshal(course.Bonuses)
	if err != nil {
		zap.L().Error("Failed to marshal course bonuses", zap.Error(err))
		return err
	}

	query := `
		INSERT INTO courses (title, description, is_public, price, author_id, is_active, cover_image_url, contraindications, recommendations, target_audience, course_basis, class_basis, bonuses)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	err = r.db.QueryRowContext(ctx, query,
		course.Title, course.Description, course.IsPublic, course.Price,
		course.AuthorID, course.IsActive, course.CoverImageURL,
		course.Contraindications, course.Recommendations, course.TargetAudience,
		course.CourseBasis, course.ClassBasis, bonusesJSON,
	).Scan(&course.ID)

	if err != nil {
		zap.L().Error("Failed to save course", zap.String("title", course.Title), zap.Error(err))
		return err
	}

	zap.L().Info("Course saved successfully", zap.Int("id", course.ID))
	return nil
}

func (r *CourseRepo) Update(ctx context.Context, course *domain.Course) error {
	bonusesJSON, err := json.Marshal(course.Bonuses)
	if err != nil {
		return err
	}

	query := `
		UPDATE courses
		SET title = $1, description = $2, is_public = $3, price = $4, is_active = $5,
		    cover_image_url = NULLIF($6, ''), 
		    contraindications = $7, 
		    recommendations = $8,
		    target_audience = $9,
		    course_basis = $10,
		    class_basis = $11,
		    bonuses = $12
		WHERE id = $13
	`

	res, err := r.db.ExecContext(ctx, query,
		course.Title, course.Description, course.IsPublic, course.Price, course.IsActive,
		course.CoverImageURL, course.Contraindications, course.Recommendations,
		course.TargetAudience, course.CourseBasis, course.ClassBasis, bonusesJSON, course.ID,
	)

	if err != nil {
		zap.L().Error("Failed to update course", zap.Int("id", course.ID), zap.Error(err))
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		zap.L().Warn("Could not get rows affected after update", zap.Int("id", course.ID), zap.Error(err))
	} else if rowsAffected == 0 {
		zap.L().Warn("Update affected 0 rows, course may not exist", zap.Int("id", course.ID))
		return domain.ErrCourseNotFound
	}

	zap.L().Debug("Course updated", zap.Int("id", course.ID))
	return nil
}

func (r *CourseRepo) UpdateStatus(ctx context.Context, id int, isActive bool) error {
	query := `UPDATE courses SET is_active = $1 WHERE id = $2`

	res, err := r.db.ExecContext(ctx, query, isActive, id)
	if err != nil {
		zap.L().Error("Failed to update course status", zap.Int("id", id), zap.Error(err))
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		zap.L().Warn("Could not get rows affected", zap.Int("id", id), zap.Error(err))
	} else if rowsAffected == 0 {
		zap.L().Warn("No course found for status update", zap.Int("id", id))
		return domain.ErrCourseNotFound
	}

	zap.L().Info("Course status updated", zap.Int("id", id), zap.Bool("is_active", isActive))
	return nil
}

func (r *CourseRepo) SetInactive(ctx context.Context, courseID int) error {
	// Можно использовать UpdateStatus, но оставим как отдельный метод для явности
	return r.UpdateStatus(ctx, courseID, false)
}
