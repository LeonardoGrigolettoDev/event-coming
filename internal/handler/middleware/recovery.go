package middleware

import (
	"net/http"

	"event-coming/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery recovers from panics and logs them
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("request_id")

				logger.Error("panic recovered",
					zap.String("request_id", requestID),
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)

				response.Error(c, http.StatusInternalServerError, "internal_error", "Internal server error")
				c.Abort()
			}
		}()

		c.Next()
	}
}
