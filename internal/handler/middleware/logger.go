package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a middleware that logs HTTP requests
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request ID if exists
		requestID := c.GetString("request_id")

		logger.Info("incoming request",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
		)

		c.Next()

		logger.Info("request completed",
			zap.String("request_id", requestID),
			zap.Int("status", c.Writer.Status()),
		)
	}
}
