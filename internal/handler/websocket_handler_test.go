package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"event-coming/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewWebSocketHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	hub := &websocket.Hub{}
	pubsub := &websocket.PubSub{}

	handler := NewWebSocketHandler(hub, pubsub, logger)
	assert.NotNil(t, handler)
}

func TestWebSocketHandler_HandleConnection_MissingOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	hub := &websocket.Hub{}
	pubsub := &websocket.PubSub{}

	handler := NewWebSocketHandler(hub, pubsub, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ws//event123", nil)
	c.Params = gin.Params{
		{Key: "organization", Value: ""},
		{Key: "event", Value: "event123"},
	}

	handler.HandleConnection(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWebSocketHandler_HandleConnection_MissingEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	hub := &websocket.Hub{}
	pubsub := &websocket.PubSub{}

	handler := NewWebSocketHandler(hub, pubsub, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ws/org123/", nil)
	c.Params = gin.Params{
		{Key: "organization", Value: "org123"},
		{Key: "event", Value: ""},
	}

	handler.HandleConnection(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
