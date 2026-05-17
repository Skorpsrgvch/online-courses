package userprofile

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	UserID int
}

type Output struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Интерфейс должен использовать GetUserByID
type UserReader interface {
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
}

type Usecase struct {
	userReader UserReader
}

func NewUsecase(userReader UserReader) (*Usecase, error) {
	if userReader == nil {
		return nil, errors.New("userReader is required")
	}
	return &Usecase{userReader: userReader}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	if input.UserID <= 0 {
		return nil, errors.New("некорректный ID пользователя")
	}

	user, err := u.userReader.GetUserByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &Output{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}
