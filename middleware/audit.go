package middleware

import (
	"fmt"
	"go-tiket-konser/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuditLogMiddleware(auditLogService service.AuditLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process the request first
		c.Next()

		method := c.Request.Method
		path := c.Request.URL.Path

		// Only log mutating actions (POST, PUT, DELETE) or specific auth requests (like verification)
		// and skip SSE streaming connections
		if method == http.MethodGet && !strings.Contains(path, "verify-otp") {
			return
		}
		if strings.Contains(path, "/stream") {
			return
		}

		statusCode := c.Writer.Status()

		// Get authenticated user info
		var userIDStr string
		var emailStr string
		var roleStr string

		if idVal, ok := c.Get("user_id"); ok {
			if str, ok := idVal.(string); ok {
				userIDStr = str
			} else if u, ok := idVal.(uuid.UUID); ok {
				userIDStr = u.String()
			}
		}

		if emailVal, ok := c.Get("email"); ok {
			if str, okStr := emailVal.(string); okStr {
				emailStr = str
			}
		}

		if roleVal, ok := c.Get("role"); ok {
			if str, okStr := roleVal.(string); okStr {
				roleStr = str
			}
		}

		// Determine action description
		action := fmt.Sprintf("%s %s", method, path)
		if strings.Contains(path, "/login") {
			action = "user_login"
		} else if strings.Contains(path, "/register") {
			action = "user_register"
		} else if strings.Contains(path, "/logout") {
			action = "user_logout"
		} else if strings.Contains(path, "/verify-otp") {
			action = "verify_otp"
		} else if strings.Contains(path, "/bookings") {
			if method == http.MethodPost {
				action = "create_booking"
			}
		} else if strings.Contains(path, "/concerts") {
			if method == http.MethodPost {
				action = "create_concert"
			} else if method == http.MethodPut {
				action = "update_concert"
			} else if method == http.MethodDelete {
				action = "delete_concert"
			}
		} else if strings.Contains(path, "/ticket-categories") {
			if method == http.MethodPost {
				action = "create_ticket_category"
			} else if method == http.MethodPut {
				action = "update_ticket_category"
			} else if method == http.MethodDelete {
				action = "delete_ticket_category"
			}
		}

		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Log it asynchronously
		go func() {
			_ = auditLogService.Log(
				action,
				method,
				path,
				statusCode,
				userIDStr,
				emailStr,
				roleStr,
				ipAddress,
				userAgent,
				"",
			)
		}()
	}
}
