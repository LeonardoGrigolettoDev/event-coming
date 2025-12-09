package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"event-coming/internal/config"
	"event-coming/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewWebhookHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	cfg := &config.WhatsAppConfig{
		WebhookVerifyToken: "test-token",
	}

	handler := NewWebhookHandler(cfg, &service.ParticipantService{}, &service.LocationService{}, logger)
	assert.NotNil(t, handler)
}

func TestWebhookHandler_VerifyWebhook_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	cfg := &config.WhatsAppConfig{
		WebhookVerifyToken: "test-verify-token",
	}

	handler := &WebhookHandler{
		cfg:    cfg,
		logger: logger,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/webhook/whatsapp?hub.mode=subscribe&hub.verify_token=test-verify-token&hub.challenge=challenge123", nil)

	handler.VerifyWebhook(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "challenge123", w.Body.String())
}

func TestWebhookHandler_VerifyWebhook_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	cfg := &config.WhatsAppConfig{
		WebhookVerifyToken: "correct-token",
	}

	handler := &WebhookHandler{
		cfg:    cfg,
		logger: logger,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/webhook/whatsapp?hub.mode=subscribe&hub.verify_token=wrong-token&hub.challenge=challenge123", nil)

	handler.VerifyWebhook(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestWebhookHandler_VerifyWebhook_WrongMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	cfg := &config.WhatsAppConfig{
		WebhookVerifyToken: "test-verify-token",
	}

	handler := &WebhookHandler{
		cfg:    cfg,
		logger: logger,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/webhook/whatsapp?hub.mode=unsubscribe&hub.verify_token=test-verify-token&hub.challenge=challenge123", nil)

	handler.VerifyWebhook(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
