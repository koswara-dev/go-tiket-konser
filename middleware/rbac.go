package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// middleware otorisasi peran
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak: role tidak ditemukan dalam token.",
			})
			return
		}

		role := roleVal.(string)
		isAllowed := false

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak: role tidak sesuai.",
			})
			return
		}

		// jika role sesuai akan dilanjutkan ke request berikutnya
		c.Next()
	}
}
