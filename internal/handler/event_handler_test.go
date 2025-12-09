package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"event-coming/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestEventHandler_Create_NoOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	// No organization_id set

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_Create_InvalidOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("organization_id", "invalid-uuid")

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_Create_NoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("organization_id", uuid.New().String())
	// No user_id set

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Create(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestEventHandler_Create_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("organization_id", uuid.New().String())
	c.Set("user_id", "invalid-uuid")

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("organization_id", uuid.New().String())
	c.Set("user_id", uuid.New().String())

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_GetByID_NoOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/events/"+uuid.New().String(), nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.GetByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_GetByID_InvalidOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/events/"+uuid.New().String(), nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	c.Set("organization_id", "invalid-uuid")

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.GetByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_GetByID_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/events/invalid-uuid", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.GetByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_Update_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/events/invalid", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_Update_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/events/"+uuid.New().String(), bytes.NewReader([]byte("invalid")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	c.Set("organization_id", uuid.New().String())

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_Delete_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/events/invalid", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_List_NoOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/events", nil)

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.List(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNewEventHandler(t *testing.T) {
	logger := zap.NewNop()
	handler := NewEventHandler(&service.EventService{}, logger)
	assert.NotNil(t, handler)
}

func TestEventHandler_Activate_NoOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/"+uuid.New().String()+"/activate", nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Activate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_Activate_InvalidOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/"+uuid.New().String()+"/activate", nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	c.Set("organization_id", "invalid-uuid")

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Activate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_Activate_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/invalid/activate", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Activate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_Cancel_NoOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/"+uuid.New().String()+"/cancel", nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Cancel(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_Cancel_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/invalid/cancel", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Cancel(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_Complete_NoOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/"+uuid.New().String()+"/complete", nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Complete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_Complete_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/invalid/complete", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.Complete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventHandler_List_InvalidOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/events", nil)
	c.Set("organization_id", "invalid-uuid")

	handler := &EventHandler{
		service: nil,
		logger:  logger,
	}

	handler.List(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
