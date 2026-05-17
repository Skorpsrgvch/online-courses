package login

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func NewUsecase(authenticator Authenticator) (*Usecase, error) {
	if authenticator == nil {
		return nil, errors.New("authenticator is required")
	}
	return &Usecase{authenticator: authenticator}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	if input.Email == "" || input.Password == "" {
		return nil, domain.ErrInvalidCredentials
	}

	user, passwordHash, err := u.authenticator.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, domain.ErrInvalidCredentials
	}

	// Сравнение пароля
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
		// Пароль не подошел
		return nil, domain.ErrInvalidCredentials
	}

	return &Output{User: user}, nil
}
