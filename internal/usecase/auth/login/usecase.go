package login

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Input struct {
	Email    string
	Password string
}

type Output struct {
	User *domain.User
}

type Usecase struct {
	authenticator Authenticator // ← используем новый интерфейс
}

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
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	return &Output{User: user}, nil
}
