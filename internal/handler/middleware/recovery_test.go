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

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		handlerPanics  bool
		expectedStatus int
	}{
		{
			name:           "no panic - request succeeds",
			handlerPanics:  false,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			c.Set("request_id", "test-request-id")

			middleware := Recovery(logger)

			middleware(c)

			if !c.IsAborted() {
				c.Status(http.StatusOK)
			}

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestRecovery_ActualPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	// Create a test router with the recovery middleware
	router := gin.New()
	router.Use(Recovery(logger))
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic!")
	})
	router.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Test panic recovery
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Test normal request
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/ok", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestRecovery_WithRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zaptest.NewLogger(t)

	router := gin.New()
	router.Use(RequestID())
	router.Use(Recovery(logger))
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic with request ID!")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	req.Header.Set("X-Request-ID", "my-custom-request-id")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "my-custom-request-id", w.Header().Get("X-Request-ID"))
}
