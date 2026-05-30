package getfull

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

// Интервал открытия модулей в тестовом режиме (5 минут)
const testModuleInterval = 5 * time.Minute

// Срок доступа к курсу (1 год)
const accessDuration = 365 * 24 * time.Hour

type Input struct {
	CourseID int
	UserID   int
	Role     string
}

type LessonOutput struct {
	ID           int     `json:"id"`
	ModuleID     int     `json:"module_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key"`
	Order        int     `json:"order"`
	IsCompleted  bool    `json:"is_completed"`
}

type ModuleOutput struct {
	ID         int            `json:"id"`
	CourseID   int            `json:"course_id"`
	Title      string         `json:"title"`
	Order      int            `json:"order"`
	WeekNumber int            `json:"week_number"`
	IsLocked   bool           `json:"is_locked"`
	UnlockDate *time.Time     `json:"unlock_date,omitempty"`
	Lessons    []LessonOutput `json:"lessons"`
}

type Output struct {
	Course          *domain.Course `json:"course"`
	Modules         []ModuleOutput `json:"modules"`
	IsPurchased     bool           `json:"is_purchased"`
	IsAccessExpired bool           `json:"is_access_expired"`
	AccessExpiresAt *time.Time     `json:"access_expires_at"`
	DaysRemaining   int            `json:"days_remaining"`
}

type CourseReader interface {
	GetByID(ctx context.Context, id int) (*domain.Course, error)
}

type ModuleReader interface {
	GetByCourseID(ctx context.Context, courseID int) ([]*domain.Module, error)
}

type LessonReader interface {
	GetByModuleID(ctx context.Context, moduleID int) ([]*domain.Lesson, error)
}

type PurchaseChecker interface {
	HasPurchased(ctx context.Context, userID, courseID int) (bool, error)
	GetPurchaseDate(ctx context.Context, userID, courseID int) (*time.Time, error)
	EnrollFree(ctx context.Context, userID, courseID, coursePrice int) error
}

type ProgressReader interface {
	GetCourseProgress(ctx context.Context, userID, courseID int) (completed int, total int, err error)
	IsLessonCompleted(ctx context.Context, userID, lessonID int) (bool, error)
}

type Usecase struct {
	courseReader    CourseReader
	moduleReader    ModuleReader
	lessonReader    LessonReader
	purchaseChecker PurchaseChecker
	progressReader  ProgressReader
	enroller        PurchaseChecker
}

func NewUsecase(
	courseReader CourseReader,
	moduleReader ModuleReader,
	lessonReader LessonReader,
	purchaseChecker PurchaseChecker,
	pr ProgressReader,
	enroller PurchaseChecker,
) (*Usecase, error) {
	if courseReader == nil || moduleReader == nil || lessonReader == nil || purchaseChecker == nil || enroller == nil {
		return nil, errors.New("all dependencies required")
	}
	return &Usecase{
		courseReader:    courseReader,
		moduleReader:    moduleReader,
		lessonReader:    lessonReader,
		purchaseChecker: purchaseChecker,
		progressReader:  pr,
		enroller:        enroller,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Debug("GetFullCourse started", zap.Int("courseID", input.CourseID), zap.Int("userID", input.UserID))

	course, err := u.courseReader.GetByID(ctx, input.CourseID)
	if err != nil {
		zap.L().Error("Failed to get course by ID", zap.Int("courseID", input.CourseID), zap.Error(err))
		return nil, err
	}

	hasFullAccess := false
	isPurchased := false
	var purchasedAt *time.Time

	// Переменные для ответа фронтенду
	isAccessExpired := false
	var accessExpiresAt *time.Time
	daysRemaining := 0

	// 1. Проверяем покупку в БД
	if input.UserID > 0 {
		purchased, err := u.purchaseChecker.HasPurchased(ctx, input.UserID, course.ID)
		if err != nil {
			zap.L().Error("Failed to check purchase status", zap.Int("userID", input.UserID), zap.Int("courseID", course.ID), zap.Error(err))
			return nil, err
		}

		// Если покупки нет, но курс бесплатный (или публичный) -> автоматически зачисляем
		if !purchased && (course.IsPublic || course.Price == 0) {
			zap.L().Info("Auto-enrolling user to free/public course", zap.Int("userID", input.UserID), zap.Int("courseID", course.ID))
			// Игнорируем ошибку, если запись уже есть (конфликт уникальности)
			_ = u.enroller.EnrollFree(ctx, input.UserID, course.ID, course.Price)

			// Пробуем проверить покупку снова (теперь она должна быть)
			purchased = true
		}

		if purchased {
			date, err := u.purchaseChecker.GetPurchaseDate(ctx, input.UserID, course.ID)
			if err != nil {
				zap.L().Error("Failed to get purchase date", zap.Error(err))
				// Если ошибка получения даты, но покупка есть - даем доступ (fallback)
				isPurchased = true
				hasFullAccess = true
			} else {
				purchasedAt = date
				zap.L().Debug("Course purchased at", zap.Time("date", *purchasedAt))

				// Расчет даты истечения
				expiryDate := purchasedAt.Add(accessDuration)
				accessExpiresAt = &expiryDate

				// Расчет оставшихся дней
				now := time.Now()
				if now.After(expiryDate) {
					isAccessExpired = true
					daysRemaining = 0
					zap.L().Warn("Access expired for user",
						zap.Int("userID", input.UserID),
						zap.Int("courseID", course.ID),
						zap.Time("expiredAt", expiryDate))
				} else {
					// Доступ активен
					durationLeft := expiryDate.Sub(now)
					daysRemaining = int(durationLeft.Hours() / 24)

					isPurchased = true
					hasFullAccess = true
				}
			}
		}
	}

	// 2. Если нет полного доступа, проверяем права (бесплатный курс или админ)
	// Этот блок теперь сработает только если авто-зачисление не прошло или курс не бесплатный
	if !hasFullAccess {
		if course.IsPublic {
			hasFullAccess = true
			zap.L().Debug("Access granted via public course flag (fallback)", zap.Int("courseID", course.ID))
		} else if input.Role == "admin" {
			hasFullAccess = true
			zap.L().Debug("Access granted via admin role", zap.Int("userID", input.UserID))
		} else {
			if isAccessExpired {
				zap.L().Info("Access denied due to expiration",
					zap.Int("userID", input.UserID),
					zap.Int("courseID", course.ID))
			} else {
				zap.L().Debug("No full access granted", zap.Bool("isPublic", course.IsPublic), zap.String("role", input.Role))
			}
		}
	}

	dbModules, err := u.moduleReader.GetByCourseID(ctx, course.ID)
	if err != nil {
		zap.L().Error("Failed to get modules", zap.Int("courseID", course.ID), zap.Error(err))
		return nil, err
	}

	modulesOut := make([]ModuleOutput, 0, len(dbModules))
	now := time.Now()

	for _, m := range dbModules {
		isLocked := false
		var unlockTime *time.Time

		isBonus := false
		titleLower := strings.ToLower(m.Title)
		if strings.Contains(titleLower, "бонус") {
			isBonus = true
		}

		// Логика блокировки модулей работает только если доступ не истек
		if !isBonus && isPurchased && purchasedAt != nil && m.WeekNumber > 1 && !isAccessExpired {
			openTime := purchasedAt.Add(time.Duration(m.WeekNumber-1) * testModuleInterval)
			unlockTime = &openTime

			if now.Before(openTime) {
				isLocked = true
				zap.L().Debug("Module is locked by time",
					zap.Int("moduleID", m.ID),
					zap.String("title", m.Title),
					zap.Int("weekNumber", m.WeekNumber),
					zap.Time("openTime", openTime))
			}
		}

		dbLessons, err := u.lessonReader.GetByModuleID(ctx, m.ID)
		if err != nil {
			zap.L().Error("Failed to get lessons for module", zap.Int("moduleID", m.ID), zap.Error(err))
			return nil, err
		}

		lessonsOut := make([]LessonOutput, 0, len(dbLessons))

		for _, l := range dbLessons {
			lessonOut := LessonOutput{
				ID:          l.ID,
				ModuleID:    l.ModuleID,
				Title:       l.Title,
				Description: l.Description,
				Order:       l.Order,
			}

			// Если доступ истек или модуль заблокирован, скрываем контент
			if isAccessExpired || isLocked {
				lessonOut.VideoEmbedID = ""
				lessonOut.PrivateKey = nil
				lessonOut.IsCompleted = false
			} else if hasFullAccess {
				lessonOut.VideoEmbedID = l.VideoEmbedID
				lessonOut.PrivateKey = l.PrivateKey

				if input.UserID > 0 {
					isCompleted, _ := u.progressReader.IsLessonCompleted(ctx, input.UserID, l.ID)
					lessonOut.IsCompleted = isCompleted
				}
			} else {
				lessonOut.VideoEmbedID = ""
				lessonOut.PrivateKey = nil
			}

			lessonsOut = append(lessonsOut, lessonOut)
		}

		if isBonus {
			isLocked = false
			unlockTime = nil
		}

		modulesOut = append(modulesOut, ModuleOutput{
			ID:         m.ID,
			CourseID:   m.CourseID,
			Title:      m.Title,
			Order:      m.Order,
			WeekNumber: m.WeekNumber,
			IsLocked:   isLocked,
			UnlockDate: unlockTime,
			Lessons:    lessonsOut,
		})
	}

	progressPercent := 0
	// Прогресс считаем только если есть полный доступ и доступ не истек
	if input.UserID > 0 && hasFullAccess && !isAccessExpired {
		completed, total, err := u.progressReader.GetCourseProgress(ctx, input.UserID, course.ID)
		if err == nil && total > 0 {
			progressPercent = int(float64(completed) / float64(total) * 100)
		} else if err != nil {
			zap.L().Warn("Failed to calculate progress", zap.Int("userID", input.UserID), zap.Int("courseID", course.ID), zap.Error(err))
		}
	}

	course.IsPurchased = isPurchased
	course.ProgressPercent = progressPercent

	zap.L().Info("GetFullCourse completed",
		zap.Int("courseID", course.ID),
		zap.Bool("isPurchased", isPurchased),
		zap.Bool("isAccessExpired", isAccessExpired),
		zap.Int("daysRemaining", daysRemaining))

	return &Output{
		Course:          course,
		Modules:         modulesOut,
		IsPurchased:     isPurchased,
		IsAccessExpired: isAccessExpired,
		AccessExpiresAt: accessExpiresAt,
		DaysRemaining:   daysRemaining,
	}, nil
}
