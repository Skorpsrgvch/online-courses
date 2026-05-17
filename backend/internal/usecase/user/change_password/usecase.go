package change_password

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
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
	GetUserByID(ctx context.Context, id int) (*domain.User, error) // Можно добавить для проверки существования, но GetPasswordHash тоже ок
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
		// Скрываем детальную ошибку, чтобы не раскрывать существование пользователя
		return nil, errors.New("неверный текущий пароль")
	}

	// 2. Сравниваем пароли
	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(input.OldPassword)); err != nil {
		return nil, errors.New("неверный текущий пароль")
	}

	// 3. Хешируем новый пароль
	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("ошибка генерации хеша")
	}

	// 4. Сохраняем
	if err := u.userRepo.UpdatePassword(ctx, input.UserID, string(newHash)); err != nil {
		return nil, errors.New("ошибка сохранения нового пароля")
	}

	return &Output{Message: "Пароль успешно изменен"}, nil
}
