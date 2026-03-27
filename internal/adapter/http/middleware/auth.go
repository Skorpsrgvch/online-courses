package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пока простая сессия или JWT — для диплома можно использовать временный токен
		// Например: Authorization: Bearer admin_token или user_token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "Bearer admin" {
			c.Set(UserIDKey, 1) // ID админа из БД
			c.Set(RoleKey, "admin")
			c.Next()
			return
		}
		if authHeader == "Bearer user" {
			c.Set(UserIDKey, 2) // ID тестового пользователя
			c.Set(RoleKey, "user")
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
	}
}

func RequireAdmin(c *gin.Context) bool {
	role, exists := c.Get(RoleKey)
	return exists && role == "admin"
}

func GetUserID(c *gin.Context) int {
	id, _ := c.Get(UserIDKey)
	return id.(int)
}
