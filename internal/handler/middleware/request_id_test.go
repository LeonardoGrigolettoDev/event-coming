package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		existingID     string
		expectCustomID bool
	}{
		{
			name:           "generates new ID when not provided",
			existingID:     "",
			expectCustomID: false,
		},
		{
			name:           "uses provided ID",
			existingID:     "custom-request-id-123",
			expectCustomID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			if tt.existingID != "" {
				c.Request.Header.Set("X-Request-ID", tt.existingID)
			}

			middleware := RequestID()
			middleware(c)

			// Check context has request_id
			requestID, exists := c.Get("request_id")
			assert.True(t, exists)
			assert.NotEmpty(t, requestID)

			// Check response header
			responseID := w.Header().Get("X-Request-ID")
			assert.NotEmpty(t, responseID)

			if tt.expectCustomID {
				assert.Equal(t, tt.existingID, requestID)
				assert.Equal(t, tt.existingID, responseID)
			} else {
				// Should be a valid UUID
				_, err := uuid.Parse(requestID.(string))
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequestID_UUIDFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	middleware := RequestID()
	middleware(c)

	requestID := c.GetString("request_id")

	// Should be a valid UUID
	parsedID, err := uuid.Parse(requestID)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, parsedID)
}

func TestRequestID_PreservesCustomID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	customID := "my-custom-trace-id-abc123"

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/test", nil)
	c.Request.Header.Set("X-Request-ID", customID)

	middleware := RequestID()
	middleware(c)

	// Context should have the custom ID
	assert.Equal(t, customID, c.GetString("request_id"))

	// Response header should have the custom ID
	assert.Equal(t, customID, w.Header().Get("X-Request-ID"))
}
