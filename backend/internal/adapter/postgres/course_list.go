package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

// ListAll возвращает список всех активных курсов
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
		return nil, err
	}
	defer rows.Close()

	var courses []*domain.Course
	for rows.Next() {
		var c domain.Course
		var bonusesRaw []byte
		var coverURL, contraindications, recommendations, targetAudience, courseBasis, classBasis sql.NullString

		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}

		// Заполняем поля из NullString
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

		// Парсим бонусы
		if len(bonusesRaw) > 0 {
			json.Unmarshal(bonusesRaw, &c.Bonuses)
		} else {
			c.Bonuses = []domain.BonusItem{}
		}

		courses = append(courses, &c)
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

	var courses []*domain.Course
	for rows.Next() {
		c, err := scanCourseHelperRow(rows)
		if err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, rows.Err()
}

// Вспомогательная функция для сканирования строки из rows (аналогична scanCourseHelper, но для rows)
func scanCourseHelperRow(rows *sql.Rows) (*domain.Course, error) {
	var c domain.Course
	var bonusesRaw []byte
	var coverURL, contraindications, recommendations, targetAudience, courseBasis, classBasis sql.NullString

	err := rows.Scan(
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
		return nil, err
	}

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

	if len(bonusesRaw) > 0 {
		json.Unmarshal(bonusesRaw, &c.Bonuses)
	} else {
		c.Bonuses = []domain.BonusItem{}
	}

	return &c, nil
}
