package register

import (
	"context"
	"errors"
	"strings"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Input struct {
	Email    string
	FullName string
	Password string
	Role     string // "user" или "admin"
}

type Usecase struct {
	userCreator UserCreator
}

func NewUsecase(userCreator UserCreator) (*Usecase, error) {
	if userCreator == nil {
		return nil, errors.New("userCreator is required")
	}
	return &Usecase{userCreator: userCreator}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	// Валидация входных данных
	if input.Email == "" || input.FullName == "" || input.Password == "" {
		return domain.ErrInvalidCredentials
	}
	if input.Role != "user" && input.Role != "admin" {
		return errors.New("invalid role")
	}

	// Создаём пользователя (без пароля!)
	user, err := domain.NewUser(input.Email, input.FullName, input.Role)
	if err != nil {
		return err

	}

	// Хэшируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	err = u.userCreator.CreateUser(ctx, user, string(passwordHash))
	if err != nil {

		if strings.Contains(err.Error(), "already exists") {
			return err
		}
		return errors.New("failed to register user")
	}

	return nil
}
