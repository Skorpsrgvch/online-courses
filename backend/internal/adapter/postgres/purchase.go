package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type PurchaseRepo struct {
	db *sql.DB
}

func NewPurchaseRepo(db *sql.DB) *PurchaseRepo {
	return &PurchaseRepo{db: db}
}

// HasPurchased проверяет, купил ли пользователь курс
func (r *PurchaseRepo) HasPurchased(ctx context.Context, userID, courseID int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM user_purchases WHERE user_id = $1 AND course_id = $2)",
		userID, courseID,
	).Scan(&exists)
	return exists, err
}

// GetByUserAndCourse возвращает запись о покупке, если она существует
// Необходимо для использования в UseCase callback
func (r *PurchaseRepo) GetByUserAndCourse(ctx context.Context, userID, courseID int) (*domain.Purchase, error) {
	p := &domain.Purchase{}
	query := `SELECT user_id, course_id, purchased_at FROM user_purchases WHERE user_id = $1 AND course_id = $2`

	err := r.db.QueryRowContext(ctx, query, userID, courseID).Scan(
		&p.UserID,
		&p.CourseID,
		&p.PurchasedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPurchaseNotFound // Убедитесь, что эта ошибка определена в domain
		}
		return nil, err
	}

	return p, nil
}

// UpdatePurchaseDate обновляет дату покупки для продления доступа
func (r *PurchaseRepo) UpdatePurchaseDate(ctx context.Context, userID, courseID int, newDate time.Time) error {
	zap.L().Info("Updating purchase date to extend access",
		zap.Int("user_id", userID),
		zap.Int("course_id", courseID),
		zap.Time("new_date", newDate))

	_, err := r.db.ExecContext(ctx,
		"UPDATE user_purchases SET purchased_at = $3 WHERE user_id = $1 AND course_id = $2",
		userID, courseID, newDate,
	)

	if err != nil {
		zap.L().Error("Failed to update purchase date", zap.Error(err))
		return err
	}

	return nil
}

// Create добавляет запись о покупке
func (r *PurchaseRepo) Create(ctx context.Context, userID, courseID int, paymentID string) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO user_purchases (user_id, course_id, purchased_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (user_id, course_id)
        DO UPDATE SET purchased_at = EXCLUDED.purchased_at
  `, userID, courseID)
	if err != nil {
		zap.L().Error("Failed to create purchase", zap.Int("user_id", userID), zap.Int("course_id", courseID), zap.Error(err))
		return err
	}
	zap.L().Info("Purchase created", zap.Int("user_id", userID), zap.Int("course_id", courseID))
	return nil
}

// GetUserCourseIDs возвращает список ID курсов
func (r *PurchaseRepo) GetUserCourseIDs(ctx context.Context, userID int) ([]int, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT course_id FROM user_purchases WHERE user_id = $1 ORDER BY purchased_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courseIDs []int
	for rows.Next() {
		var cid int
		if err := rows.Scan(&cid); err != nil {
			return nil, err
		}
		courseIDs = append(courseIDs, cid)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courseIDs, nil
}

// GetByUser возвращает полный список покупок с данными курсов
func (r *PurchaseRepo) GetByUser(ctx context.Context, userID int) ([]domain.CourseWithDate, error) {
	query := `
		SELECT c.id, c.title, COALESCE(c.description, ''), c.price, COALESCE(c.cover_image_url, ''), up.purchased_at
		FROM user_purchases up
		JOIN courses c ON c.id = up.course_id
		WHERE up.user_id = $1
		ORDER BY up.purchased_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []domain.CourseWithDate

	for rows.Next() {
		var p domain.CourseWithDate
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.CoverImageURL, &p.PurchasedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if purchases == nil {
		return []domain.CourseWithDate{}, nil
	}

	return purchases, nil
}

// GrantAccess предоставляет доступ к курсу (идемпотентно)
func (r *PurchaseRepo) GrantAccess(ctx context.Context, userID, courseID int) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO user_purchases (user_id, course_id) 
        VALUES ($1, $2)
        ON CONFLICT (user_id, course_id) DO NOTHING
    `, userID, courseID)
	return err
}

// EnrollFree записывает пользователя на бесплатный курс
func (r *PurchaseRepo) EnrollFree(ctx context.Context, userID, courseID int, coursePrice int) error {
	zap.L().Debug("EnrollFree started", zap.Int("user_id", userID), zap.Int("course_id", courseID), zap.Int("price", coursePrice))

	if coursePrice > 0 {
		err := errors.New("cannot enroll free for a paid course")
		zap.L().Warn("EnrollFree validation failed", zap.Int("course_id", courseID), zap.Int("price", coursePrice))
		return err
	}

	query := `
    INSERT INTO user_purchases (user_id, course_id, purchased_at) 
    VALUES ($1, $2, NOW())
    ON CONFLICT (user_id, course_id)
    DO UPDATE SET purchased_at = EXCLUDED.purchased_at
    `

	res, err := r.db.ExecContext(ctx, query, userID, courseID)
	if err != nil {
		zap.L().Error("EnrollFree DB execution failed", zap.Error(err))
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		zap.L().Warn("EnrollFree failed to get rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		zap.L().Info("EnrollFree: no rows inserted (likely duplicate)", zap.Int("user_id", userID), zap.Int("course_id", courseID))
	} else {
		zap.L().Info("EnrollFree: successfully enrolled", zap.Int("user_id", userID), zap.Int("course_id", courseID))
	}

	return nil
}
