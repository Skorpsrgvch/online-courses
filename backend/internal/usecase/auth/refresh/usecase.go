package refresh

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
)

func NewUsecase(repo RefreshTokenRepo, userRepo UserRepo) (*Usecase, error) {
	if repo == nil || userRepo == nil {
		return nil, errors.New("dependencies required")
	}
	return &Usecase{repo: repo, userRepo: userRepo}, nil
}

func (u *Usecase) Execute(ctx context.Context, refreshToken string) (*Output, error) {
	// 1. Парсим Refresh токен, чтобы достать UserID и проверить подпись
	claims, err := middleware.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token signature")
	}

	// 2. Проверяем, что это именно refresh токен
	// Для простоты проверяем, не истек ли он по времени JWT
	if time.Now().UTC().After(claims.ExpiresAt.Time) {
		return nil, errors.New("refresh token expired")
	}

	// 3. Проверяем наличие хеша в БД (не был ли отзыв)
	if err := u.repo.Validate(ctx, claims.UserID, refreshToken); err != nil {
		return nil, errors.New("refresh token not found in database or revoked")
	}

	// 4. Получаем актуальные данные пользователя (вдруг сменили роль или имя)
	user, err := u.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 5. Генерируем новую пару токенов
	newAccessToken, newRefreshToken, err := middleware.GenerateToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, err
	}

	// Удаляем старый токен при ротации (повышает безопасность)
	u.repo.DeleteByUser(ctx, user.ID)

	// 6. Сохраняем новый токен в БД
	expiresAt := time.Now().UTC().Add(7 * 24 * time.Hour)
	if err := u.repo.Save(ctx, user.ID, newRefreshToken, expiresAt); err != nil {
		return nil, err
	}

	return &Output{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
