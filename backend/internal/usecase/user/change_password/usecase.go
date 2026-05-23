package change_password

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Input struct {
	UserID      int
	OldPassword string
	NewPassword string
}

type Output struct {
	Message string
}

type UserRepo interface {
	GetPasswordHash(ctx context.Context, userID int) (string, error)
	UpdatePassword(ctx context.Context, userID int, hash string) error
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
	zap.L().Debug("ChangePassword started", zap.Int("userID", input.UserID))

	if input.UserID <= 0 {
		return nil, errors.New("некорректный ID пользователя")
	}

	if input.OldPassword == "" || input.NewPassword == "" {
		return nil, errors.New("старый и новый пароль обязательны")
	}

	if len(input.NewPassword) < 6 {
		return nil, errors.New("новый пароль должен быть не менее 6 символов")
	}

	// 1. Получаем хеш старого пароля
	currentHash, err := u.userRepo.GetPasswordHash(ctx, input.UserID)
	if err != nil {
		zap.L().Warn("Failed to get password hash", zap.Int("userID", input.UserID), zap.Error(err))
		// Скрываем детальную ошибку, чтобы не раскрывать существование пользователя
		return nil, errors.New("неверный текущий пароль")
	}

	// 2. Сравниваем пароли
	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(input.OldPassword)); err != nil {
		zap.L().Warn("Invalid old password provided", zap.Int("userID", input.UserID))
		return nil, errors.New("неверный текущий пароль")
	}

	// 3. Хешируем новый пароль
	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("Failed to generate password hash", zap.Error(err))
		return nil, errors.New("ошибка генерации хеша")
	}

	// 4. Сохраняем
	if err := u.userRepo.UpdatePassword(ctx, input.UserID, string(newHash)); err != nil {
		zap.L().Error("Failed to update password in DB", zap.Int("userID", input.UserID), zap.Error(err))
		return nil, errors.New("ошибка сохранения нового пароля")
	}

	zap.L().Info("Password changed successfully", zap.Int("userID", input.UserID))
	return &Output{Message: "Пароль успешно изменен"}, nil
}
