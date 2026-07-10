package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS) configurations.
// For development (APP_ENV != "production"), it allows all domains.
// For production (APP_ENV == "production"), it only allows https://domain.id.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		env := os.Getenv("APP_ENV")

		allowOrigin := ""
		if env == "production" {
			if origin == "https://domain.id" {
				allowOrigin = "https://domain.id"
			}
		} else {
			// In development, allow all domains.
			// Reflecting request's Origin allows requests with credentials to pass successfully.
			if origin != "" {
				allowOrigin = origin
			} else {
				allowOrigin = "*"
			}
		}

		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-api-key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
