package userprofile

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
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
	zap.L().Debug("GetUserProfile started", zap.Int("userID", input.UserID))

	if input.UserID <= 0 {
		return nil, errors.New("некорректный ID пользователя")
	}

	user, err := u.userReader.GetUserByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			zap.L().Warn("User not found", zap.Int("userID", input.UserID))
		} else {
			zap.L().Error("Failed to get user by ID", zap.Int("userID", input.UserID), zap.Error(err))
		}
		return nil, err
	}

	zap.L().Info("User profile retrieved", zap.Int("userID", input.UserID), zap.String("email", user.Email))
	return &Output{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}
