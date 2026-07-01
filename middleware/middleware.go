package middleware

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// ApiKeyAuth middleware untuk validasi x-api-key
func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}
		apiKey := c.GetHeader("x-api-key")
		if apiKey != "juara-coding-super-secret" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// IPRateLimiter simple rate limiter
type IPRateLimiter struct {
	mu   sync.Mutex
	hits map[string]int
}

var limiter = IPRateLimiter{
	hits: make(map[string]int),
}

// RateLimiter limit request per user
func RateLimiter(maxRequest int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		limiter.mu.Lock()
		currentHits := limiter.hits[clientIP]

		if currentHits >= maxRequest {
			limiter.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			return
		}

		limiter.hits[clientIP] = currentHits + 1
		limiter.mu.Unlock()

		c.Next()
	}
}
