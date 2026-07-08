package middleware

import (
	"go-tiket-konser/config"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

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

// IPRateLimiter simple memory rate limiter with a reset timer
type IPRateLimiter struct {
	mu    sync.Mutex
	hits  map[string]int
	reset time.Time
}

var limiter = IPRateLimiter{
	hits:  make(map[string]int),
	reset: time.Now().Add(1 * time.Minute),
}

func (l *IPRateLimiter) checkLimit(ip string, maxRequest int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if now.After(l.reset) {
		l.hits = make(map[string]int)
		l.reset = now.Add(1 * time.Minute)
	}

	l.hits[ip]++
	return l.hits[ip] <= maxRequest
}

// RateLimiter limit request per user
func RateLimiter(maxRequest int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("APP_ENV") != "production" {
			c.Next()
			return
		}
		clientIP := c.ClientIP()

		if config.RedisClient != nil {
			ctx := c.Request.Context()
			key := "rate_limit:" + clientIP
			pipe := config.RedisClient.Pipeline()
			incr := pipe.Incr(ctx, key)
			pipe.Expire(ctx, key, 1*time.Minute)
			_, err := pipe.Exec(ctx)
			if err == nil {
				if incr.Val() > int64(maxRequest) {
					c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
					return
				}
				c.Next()
				return
			}
		}

		// Fallback to memory-efficient limiter
		if !limiter.checkLimit(clientIP, maxRequest) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			return
		}
		c.Next()
	}
}
