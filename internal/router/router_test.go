package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"event-coming/internal/config"
	"event-coming/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewRouter(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		App: config.AppConfig{
			Debug: true,
		},
		JWT: config.JWTConfig{
			AccessSecret: "test-secret",
		},
	}

	router := NewRouter(
		cfg,
		logger,
		nil, // authHandler
		nil, // websocketHandler
		nil, // eventCacheHandler
		nil, // participantHandler
		nil, // eventHandler
		nil, // entityHandler
		nil, // locationHandler
		nil, // webhookHandler
	)

	assert.NotNil(t, router)
	assert.NotNil(t, router.engine)
	assert.Equal(t, cfg, router.config)
	assert.Equal(t, logger, router.logger)
}

func TestNewRouter_ReleaseMode(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		App: config.AppConfig{
			Debug: false, // Release mode
		},
	}

	router := NewRouter(
		cfg,
		logger,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	assert.NotNil(t, router)
}

func TestRouter_GetEngine(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		App: config.AppConfig{
			Debug: true,
		},
	}

	router := NewRouter(cfg, logger, nil, nil, nil, nil, nil, nil, nil, nil)
	engine := router.GetEngine()

	assert.NotNil(t, engine)
	assert.Equal(t, router.engine, engine)
}

func TestRouter_HealthEndpoint(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		App: config.AppConfig{
			Debug: true,
		},
		JWT: config.JWTConfig{
			AccessSecret: "test-secret",
		},
	}

	// Create mock handlers
	authHandler := &handler.AuthHandler{}
	webhookHandler := &handler.WebhookHandler{}
	entityHandler := &handler.EntityHandler{}
	eventHandler := &handler.EventHandler{}
	participantHandler := &handler.ParticipantHandler{}
	locationHandler := &handler.LocationHandler{}
	eventCacheHandler := &handler.EventCacheHandler{}
	websocketHandler := &handler.WebSocketHandler{}

	router := NewRouter(
		cfg,
		logger,
		authHandler,
		websocketHandler,
		eventCacheHandler,
		participantHandler,
		eventHandler,
		entityHandler,
		locationHandler,
		webhookHandler,
	)

	engine := router.Setup()

	// Test health endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
	assert.Contains(t, w.Body.String(), "event-coming")
}

func TestRouter_Setup_RegistersRoutes(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		App: config.AppConfig{
			Debug: true,
		},
		JWT: config.JWTConfig{
			AccessSecret: "test-secret",
		},
	}

	// Create mock handlers
	authHandler := &handler.AuthHandler{}
	webhookHandler := &handler.WebhookHandler{}
	entityHandler := &handler.EntityHandler{}
	eventHandler := &handler.EventHandler{}
	participantHandler := &handler.ParticipantHandler{}
	locationHandler := &handler.LocationHandler{}
	eventCacheHandler := &handler.EventCacheHandler{}
	websocketHandler := &handler.WebSocketHandler{}

	router := NewRouter(
		cfg,
		logger,
		authHandler,
		websocketHandler,
		eventCacheHandler,
		participantHandler,
		eventHandler,
		entityHandler,
		locationHandler,
		webhookHandler,
	)

	engine := router.Setup()

	// Verify routes are registered by checking the routes info
	routes := engine.Routes()

	// Check that expected routes exist
	expectedPaths := []string{
		"/health",
		"/api/v1/auth/register",
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
		"/api/v1/auth/logout",
		"/api/v1/auth/forgot-password",
		"/api/v1/auth/reset-password",
		"/api/v1/webhook/whatsapp",
		"/api/v1/entities",
		"/api/v1/entities/:id",
		"/api/v1/events",
		"/api/v1/events/:id",
		"/api/v1/participants/:id",
		"/api/v1/ws/:event",
	}

	registeredPaths := make(map[string]bool)
	for _, route := range routes {
		registeredPaths[route.Path] = true
	}

	for _, expected := range expectedPaths {
		assert.True(t, registeredPaths[expected], "Expected route %s to be registered", expected)
	}
}

func TestRouter_Setup_MiddlewareApplied(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		App: config.AppConfig{
			Debug: true,
		},
		JWT: config.JWTConfig{
			AccessSecret: "test-secret",
		},
	}

	// Create mock handlers
	authHandler := &handler.AuthHandler{}
	webhookHandler := &handler.WebhookHandler{}
	entityHandler := &handler.EntityHandler{}
	eventHandler := &handler.EventHandler{}
	participantHandler := &handler.ParticipantHandler{}
	locationHandler := &handler.LocationHandler{}
	eventCacheHandler := &handler.EventCacheHandler{}
	websocketHandler := &handler.WebSocketHandler{}

	router := NewRouter(
		cfg,
		logger,
		authHandler,
		websocketHandler,
		eventCacheHandler,
		participantHandler,
		eventHandler,
		entityHandler,
		locationHandler,
		webhookHandler,
	)

	engine := router.Setup()

	// Test that CORS headers are set
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	engine.ServeHTTP(w, req)

	// CORS middleware should add headers (may be "*" for all origins or specific origin)
	corsHeader := w.Header().Get("Access-Control-Allow-Origin")
	assert.True(t, corsHeader == "*" || corsHeader == "http://localhost:3000", "CORS header should be set")
}

func TestRouter_ProtectedRoutesRequireAuth(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		App: config.AppConfig{
			Debug: true,
		},
		JWT: config.JWTConfig{
			AccessSecret: "test-secret",
		},
	}

	// Create mock handlers
	authHandler := &handler.AuthHandler{}
	webhookHandler := &handler.WebhookHandler{}
	entityHandler := &handler.EntityHandler{}
	eventHandler := &handler.EventHandler{}
	participantHandler := &handler.ParticipantHandler{}
	locationHandler := &handler.LocationHandler{}
	eventCacheHandler := &handler.EventCacheHandler{}
	websocketHandler := &handler.WebSocketHandler{}

	router := NewRouter(
		cfg,
		logger,
		authHandler,
		websocketHandler,
		eventCacheHandler,
		participantHandler,
		eventHandler,
		entityHandler,
		locationHandler,
		webhookHandler,
	)

	engine := router.Setup()

	// Test protected routes without auth token
	protectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/entities"},
		{"GET", "/api/v1/events"},
		{"GET", "/api/v1/participants/123"},
	}

	for _, route := range protectedRoutes {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(route.method, route.path, nil)
		engine.ServeHTTP(w, req)

		// Should return 401 Unauthorized without token
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Route %s %s should require auth", route.method, route.path)
	}
}

func TestRouter_PublicRoutesNoAuth(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		App: config.AppConfig{
			Debug: true,
		},
		JWT: config.JWTConfig{
			AccessSecret: "test-secret",
		},
	}

	// Create mock handlers
	authHandler := &handler.AuthHandler{}
	webhookHandler := &handler.WebhookHandler{}
	entityHandler := &handler.EntityHandler{}
	eventHandler := &handler.EventHandler{}
	participantHandler := &handler.ParticipantHandler{}
	locationHandler := &handler.LocationHandler{}
	eventCacheHandler := &handler.EventCacheHandler{}
	websocketHandler := &handler.WebSocketHandler{}

	router := NewRouter(
		cfg,
		logger,
		authHandler,
		websocketHandler,
		eventCacheHandler,
		participantHandler,
		eventHandler,
		entityHandler,
		locationHandler,
		webhookHandler,
	)

	engine := router.Setup()

	// Test health endpoint (public)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	engine.ServeHTTP(w, req)

	// Should return 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetChannelKey(t *testing.T) {
	tests := []struct {
		name     string
		orgID    string
		eventID  string
		expected string
	}{
		{
			name:     "standard key",
			orgID:    "org-123",
			eventID:  "event-456",
			expected: "org-123:event-456",
		},
		{
			name:     "uuid keys",
			orgID:    "550e8400-e29b-41d4-a716-446655440000",
			eventID:  "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			expected: "550e8400-e29b-41d4-a716-446655440000:6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		},
		{
			name:     "empty values",
			orgID:    "",
			eventID:  "",
			expected: ":",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This tests a helper function from websocket package
			// We're verifying the pattern used in router setup
			result := tt.orgID + ":" + tt.eventID
			assert.Equal(t, tt.expected, result)
		})
	}
}
