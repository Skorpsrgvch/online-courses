package middleware

import (
	"strings"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
	EmailKey  = "email"
	NameKey   = "name"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// В продакшене лучше паниковать, если секрет не задан, но для совместимости оставим warning
		zap.L().Warn("JWT_SECRET is not set, using empty secret (unsafe for production)")
	}
	jwtSecret = []byte(secret)
}

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int, email, name, role string) (accessToken, refreshToken string, err error) {
	now := time.Now().UTC()

	accessClaims := Claims{
		UserID: userID,
		Email:  email,
		Name:   name,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "online-courses",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = token.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "online-courses",
		},
	}

	refreshTokenStr := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenStr.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// AuthMiddleware проверяет токен. Если токен есть, но невалиден — возвращает 401.
// Если токена нет — пропускает запрос дальше (для публичных эндпоинтов внутри защищенной группы).
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			zap.L().Debug("Invalid authorization header format", zap.String("header", authHeader))
			c.JSON(401, gin.H{"error": "неверный формат токена"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ParseToken(tokenString)

		if err != nil {
			zap.L().Info("Invalid token", zap.Error(err))
			c.JSON(401, gin.H{"error": "неавторизован"})
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(RoleKey, claims.Role)
		c.Set(EmailKey, claims.Email)
		c.Set(NameKey, claims.Name)

		c.Next()
	}
}

// RequireAuth гарантирует, что пользователь авторизован.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			zap.L().Debug("Access denied: user not authenticated")
			c.JSON(401, gin.H{"error": "требуется авторизация"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAdmin проверяет роль администратора. Возвращает false, если роль не admin.
func RequireAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			zap.L().Debug("Access denied: user not authenticated")
			c.JSON(401, gin.H{"error": "требуется авторизация"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			zap.L().Debug("Admin access denied: user not authenticated")
			c.JSON(401, gin.H{"error": "требуется авторизация"})
			c.Abort()
			return
		}

		role, exists := c.Get(RoleKey)
		if !exists || role != "admin" {
			zap.L().Warn("Admin access denied", zap.Any("role", role))
			c.JSON(403, gin.H{"error": "доступ закрыт"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireAdmin(c *gin.Context) bool {
	role, exists := c.Get(RoleKey)
	if !exists {
		return false
	}
	isAdmin := role == "admin"
	if !isAdmin {
		zap.L().Warn("Access denied: user is not admin", zap.Any("role", role))
	}
	return isAdmin
}

func GetUserID(c *gin.Context) int {
	val, exists := c.Get(UserIDKey)
	if !exists {
		return 0
	}
	userID, ok := val.(int)
	if !ok {
		return 0
	}
	return userID
}
