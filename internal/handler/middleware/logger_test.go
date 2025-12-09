package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name      string
		requestID string
		method    string
		path      string
	}{
		{
			name:      "with request ID",
			requestID: "test-request-id-123",
			method:    http.MethodGet,
			path:      "/api/v1/test",
		},
		{
			name:      "without request ID",
			requestID: "",
			method:    http.MethodPost,
			path:      "/api/v1/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(tt.method, tt.path, nil)

			if tt.requestID != "" {
				c.Set("request_id", tt.requestID)
			}

			var handlerCalled bool
			middleware := Logger(logger)
			middleware(c)

			if !c.IsAborted() {
				handlerCalled = true
				c.Status(http.StatusOK)
			}

			assert.True(t, handlerCalled, "handler should be called")
		})
	}
}

func TestLogger_WithZapNop(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Set("request_id", "nop-test-id")

	middleware := Logger(logger)
	middleware(c)

	// Should not panic with nop logger
	assert.False(t, c.IsAborted())
}
