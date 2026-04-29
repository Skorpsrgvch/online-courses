package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type CourseRepo struct {
	db *sql.DB
}

func NewCourseRepo(db *sql.DB) *CourseRepo {
	return &CourseRepo{db: db}
}

// scanCourseHelper помогает сканировать строку результата в структуру Course
// Принимает указатель на row и заполняет курс, включая новые поля
func scanCourseHelper(row *sql.Row) (*domain.Course, error) {
	var c domain.Course
	var bonusesRaw []byte
	var coverURL, contraindications, recommendations, targetAudience, courseBasis, classBasis sql.NullString

	err := row.Scan(
		&c.ID,
		&c.Title,
		&c.Description,
		&c.IsPublic,
		&c.Price,
		&c.AuthorID,
		&c.IsActive,
		&coverURL,
		&contraindications,
		&recommendations,
		&targetAudience,
		&courseBasis,
		&classBasis,
		&bonusesRaw,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCourseNotFound
		}
		return nil, err
	}

	// Обработка NULL значений для строк
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

	// Парсим JSON бонусов
	if len(bonusesRaw) > 0 {
		if err := json.Unmarshal(bonusesRaw, &c.Bonuses); err != nil {
			c.Bonuses = []domain.BonusItem{}
		}
	} else {
		c.Bonuses = []domain.BonusItem{}
	}

	return &c, nil
}

// GetByID возвращает курс по ID
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
		WHERE id = $1 AND is_active = true
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanCourseHelper(row)
}

// Save создаёт новый курс
func (r *CourseRepo) Save(ctx context.Context, course *domain.Course) error {
	query := `
		INSERT INTO courses (title, description, is_public, price, author_id, is_active, cover_image_url, contraindications, recommendations, target_audience, course_basis, class_basis, bonuses)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	bonusesJSON, err := json.Marshal(course.Bonuses)
	if err != nil {
		return err
	}

	return r.db.QueryRowContext(ctx, query,
		course.Title,
		course.Description,
		course.IsPublic,
		course.Price,
		course.AuthorID,
		course.IsActive,
		course.CoverImageURL,
		course.Contraindications,
		course.Recommendations,
		course.TargetAudience,
		course.CourseBasis,
		course.ClassBasis,
		bonusesJSON,
	).Scan(&course.ID)
}

// Update обновляет курс
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
	_, err = r.db.ExecContext(ctx, query,
		course.Title,
		course.Description,
		course.IsPublic,
		course.Price,
		course.IsActive,
		course.CoverImageURL,
		course.Contraindications,
		course.Recommendations,
		course.TargetAudience,
		course.CourseBasis,
		course.ClassBasis,
		bonusesJSON,
		course.ID,
	)
	return err
}

func (r *CourseRepo) SetInactive(ctx context.Context, courseID int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE courses SET is_active = false WHERE id = $1`,
		courseID,
	)
	return err
}
