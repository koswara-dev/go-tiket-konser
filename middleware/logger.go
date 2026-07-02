package middleware

import (
	"go-tiket-konser/utils/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		logFields := logrus.Fields{
			"status":     statusCode,
			"latency_ms": latency.Milliseconds(),
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
		}

		// Jika status HTTP >= 500 (Kesalahan Server), catat dengan level Error
		if statusCode >= 500 {
			logger.Log.WithFields(logFields).Error("Terjadi kesalahan server internal pada route ini")
		} else if statusCode >= 400 {
			// Jika status HTTP 4xx (Kesalahan Input Klien), catat dengan level Warn
			logger.Log.WithFields(logFields).Warn("Permintaan klien tidak valid atau otorisasi gagal")
		} else {
			// Jika sukses 2xx / 3xx, catat dengan level Info
			logger.Log.WithFields(logFields).Info("HTTP Request sukses diproses")
		}
	}
}
