package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"event-coming/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func generateValidToken(secret string, userID, orgID uuid.UUID) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":         userID.String(),
		"email":           "test@example.com",
		"organization_id": orgID.String(),
		"role":            "org_admin",
		"exp":             time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func generateExpiredToken(secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": uuid.New().String(),
		"exp":     time.Now().Add(-time.Hour).Unix(), // Expired
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestAuthMiddleware(t *testing.T) {
	jwtConfig := &config.JWTConfig{
		AccessSecret: "test-secret",
	}

	tests := []struct {
		name           string
		setupRequest   func(*http.Request)
		expectedStatus int
		checkContext   func(*testing.T, *gin.Context)
	}{
		{
			name: "valid token",
			setupRequest: func(r *http.Request) {
				userID := uuid.New()
				orgID := uuid.New()
				token := generateValidToken(jwtConfig.AccessSecret, userID, orgID)
				r.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				_, exists := c.Get("user_id")
				assert.True(t, exists)
			},
		},
		{
			name: "missing authorization header",
			setupRequest: func(r *http.Request) {
				// No header
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid authorization format - no Bearer",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "InvalidFormat token123")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid authorization format - single part",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "token123")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "expired token",
			setupRequest: func(r *http.Request) {
				token := generateExpiredToken(jwtConfig.AccessSecret)
				r.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid token signature",
			setupRequest: func(r *http.Request) {
				token := generateValidToken("wrong-secret", uuid.New(), uuid.New())
				r.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "malformed token",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid-token-format")
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			tt.setupRequest(c.Request)

			var handlerCalled bool
			middleware := AuthMiddleware(jwtConfig)

			// Create a test handler that sets handlerCalled to true
			c.Set("_handler_called", false)
			middleware(c)

			if !c.IsAborted() {
				// Next handler would be called
				handlerCalled = true
				c.Status(http.StatusOK)
			}

			if tt.expectedStatus == http.StatusOK {
				assert.True(t, handlerCalled, "handler should be called")
				if tt.checkContext != nil {
					tt.checkContext(t, c)
				}
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAuthMiddleware_ClaimsExtraction(t *testing.T) {
	jwtConfig := &config.JWTConfig{
		AccessSecret: "test-secret",
	}

	userID := uuid.New()
	orgID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	token := generateValidToken(jwtConfig.AccessSecret, userID, orgID)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	middleware := AuthMiddleware(jwtConfig)
	middleware(c)

	assert.False(t, c.IsAborted())

	// Check extracted claims
	extractedUserID, exists := c.Get("user_id")
	assert.True(t, exists)
	assert.Equal(t, userID, extractedUserID)

	extractedEmail, exists := c.Get("email")
	assert.True(t, exists)
	assert.Equal(t, "test@example.com", extractedEmail)

	extractedOrgID, exists := c.Get("organization_id")
	assert.True(t, exists)
	assert.Equal(t, orgID, extractedOrgID)

	extractedRole, exists := c.Get("role")
	assert.True(t, exists)
	// Role is stored as domain.UserRole, not string
	assert.NotEmpty(t, extractedRole)
}
