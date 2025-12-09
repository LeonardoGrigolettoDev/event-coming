package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkHeaders   bool
	}{
		{
			name:           "regular GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "regular POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "OPTIONS preflight request",
			method:         http.MethodOptions,
			expectedStatus: 204,
			checkHeaders:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(tt.method, "/", nil)

			var handlerCalled bool
			middleware := CORS()
			middleware(c)

			if !c.IsAborted() {
				handlerCalled = true
				c.Status(http.StatusOK)
			}

			if tt.checkHeaders {
				assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
				assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
				assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
			}

			if tt.method == http.MethodOptions {
				assert.False(t, handlerCalled, "handler should not be called for OPTIONS")
				assert.Equal(t, 204, w.Code)
			} else {
				assert.True(t, handlerCalled, "handler should be called")
			}
		})
	}
}
