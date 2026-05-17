package update_profile

import (
	"context"
	"errors"
	"strings"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	UserID   int
	NewName  *string
	NewEmail *string
}

type Output struct {
	Message string
}

type UserRepo interface {
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	UpdateName(ctx context.Context, userID int, name string) error
	UpdateEmail(ctx context.Context, userID int, email string) error
	GetUserByEmailForCheck(ctx context.Context, email string) (*domain.User, error)
}

type Usecase struct {
	userRepo UserRepo
}

func NewUsecase(userRepo UserRepo) (*Usecase, error) {
	if userRepo == nil {
		return nil, errors.New("user repo is required")
	}
	return &Usecase{userRepo: userRepo}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	if input.UserID <= 0 {
		return nil, errors.New("некорректный ID пользователя")
	}

	// 1. Проверяем существование пользователя
	user, err := u.userRepo.GetUserByID(ctx, input.UserID) // Изменено на GetUserByID
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	updated := false

	// 2. Обновление имени
	if input.NewName != nil {
		cleanName := strings.TrimSpace(*input.NewName)
		if cleanName == "" {
			return nil, errors.New("имя не может быть пустым")
		}
		if cleanName != user.Name {
			if err := u.userRepo.UpdateName(ctx, input.UserID, cleanName); err != nil {
				return nil, errors.New("ошибка при обновлении имени")
			}
			updated = true
		}
	}

	// 3. Обновление Email
	if input.NewEmail != nil {
		cleanEmail := strings.TrimSpace(*input.NewEmail)
		if cleanEmail == "" {
			return nil, errors.New("email не может быть пустым")
		}
		if cleanEmail != user.Email {
			// Проверяем, не занят ли email другим пользователем
			existingUser, err := u.userRepo.GetUserByEmailForCheck(ctx, cleanEmail)
			if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
				return nil, errors.New("ошибка проверки email")
			}

			// Если пользователь найден И это не текущий пользователь (на всякий случай)
			if existingUser != nil && existingUser.ID != input.UserID {
				return nil, errors.New("этот email уже занят другим пользователем")
			}

			// Пытаемся обновить
			if err := u.userRepo.UpdateEmail(ctx, input.UserID, cleanEmail); err != nil {
				// Если репозиторий вернул ошибку про дубликат, пробрасываем её как BadRequest
				if strings.Contains(err.Error(), "already exists") {
					return nil, common.BadRequestError("этот email уже зарегистрирован")
				}
				return nil, errors.New("ошибка при обновлении email")
			}
			updated = true
		}
	}

	if !updated {
		return &Output{Message: "Нет данных для обновления"}, nil
	}

	return &Output{Message: "Профиль успешно обновлен"}, nil
}
