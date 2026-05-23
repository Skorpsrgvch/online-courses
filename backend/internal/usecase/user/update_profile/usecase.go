package update_profile

import (
	"context"
	"errors"
	"strings"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
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
	zap.L().Debug("UpdateProfile started", zap.Int("userID", input.UserID))

	if input.UserID <= 0 {
		return nil, errors.New("некорректный ID пользователя")
	}

	// 1. Проверяем существование пользователя
	user, err := u.userRepo.GetUserByID(ctx, input.UserID)
	if err != nil {
		zap.L().Warn("User not found for update", zap.Int("userID", input.UserID), zap.Error(err))
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
			zap.L().Debug("Updating user name", zap.String("old", user.Name), zap.String("new", cleanName))
			if err := u.userRepo.UpdateName(ctx, input.UserID, cleanName); err != nil {
				zap.L().Error("Failed to update name", zap.Error(err))
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
				zap.L().Error("Error checking email availability", zap.Error(err))
				return nil, errors.New("ошибка проверки email")
			}

			// Если пользователь найден И это не текущий пользователь
			if existingUser != nil && existingUser.ID != input.UserID {
				zap.L().Info("Email already taken", zap.String("email", cleanEmail))
				return nil, errors.New("этот email уже занят другим пользователем")
			}

			// Пытаемся обновить
			zap.L().Debug("Updating user email", zap.String("old", user.Email), zap.String("new", cleanEmail))
			if err := u.userRepo.UpdateEmail(ctx, input.UserID, cleanEmail); err != nil {
				if strings.Contains(err.Error(), "already exists") {
					return nil, common.BadRequestError("этот email уже зарегистрирован")
				}
				zap.L().Error("Failed to update email", zap.Error(err))
				return nil, errors.New("ошибка при обновлении email")
			}
			updated = true
		}
	}

	if !updated {
		zap.L().Info("No changes detected in profile update", zap.Int("userID", input.UserID))
		return &Output{Message: "Нет данных для обновления"}, nil
	}

	zap.L().Info("Profile updated successfully", zap.Int("userID", input.UserID))
	return &Output{Message: "Профиль успешно обновлен"}, nil
}
