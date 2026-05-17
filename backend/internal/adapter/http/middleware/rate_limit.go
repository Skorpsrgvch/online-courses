package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Visitor struct {
	limiter *rate.Limiter
}

func RateLimitMiddleware(limit rate.Limit, burst int) gin.HandlerFunc {
	var (
		visitors = make(map[string]*Visitor)
		mu       sync.Mutex
	)

	getLimiter := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()

		v, exists := visitors[ip]
		if !exists {
			limiter := rate.NewLimiter(limit, burst)
			visitors[ip] = &Visitor{limiter: limiter}
			return limiter
		}
		return v.limiter
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Слишком много запросов. Попробуйте позже.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
