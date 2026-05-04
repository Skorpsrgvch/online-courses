package middleware

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		secret = "online-courses-super-secret-jwt-key-2026"
	}
	jwtSecret = []byte(secret)
}

// Claims — данные, хранимые в JWT
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken создаёт JWT + Refresh-токен
func GenerateToken(userID int, email, name, role string) (accessToken, refreshToken string, err error) {
	now := time.Now().UTC()

	// Access token — 15 минут
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

	// Refresh token — 30 дней
	refreshClaims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * 24 * time.Hour)),
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

// ParseToken валидирует JWT и возвращает Claims
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ParseToken(tokenString)

		if err != nil {
			log.Printf("Invalid token: %v", err)
			c.Next()
			return
		}

		// Если токен валиден — кладем данные в контекст
		c.Set(UserIDKey, claims.UserID)
		c.Set(RoleKey, claims.Role)
		c.Set(EmailKey, claims.Email)
		c.Set(NameKey, claims.Name)

		c.Next()
	}
}

func RequireAdmin(c *gin.Context) bool {
	role, exists := c.Get(RoleKey)
	if !exists {
		return false
	}
	return role == "admin"
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
