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

func TestParticipantHandler_Create_NoOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/"+uuid.New().String()+"/participants", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_Create_InvalidOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/"+uuid.New().String()+"/participants", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	c.Set("organization_id", "invalid-uuid")

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_Create_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/invalid/participants", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/"+uuid.New().String()+"/participants", bytes.NewReader([]byte("invalid")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	c.Set("organization_id", uuid.New().String())

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_GetByID_InvalidParticipantID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/participants/invalid", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.GetByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_Update_InvalidParticipantID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/participants/invalid", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_Delete_InvalidParticipantID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/participants/invalid", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_ListByEvent_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/events/invalid/participants", nil)
	c.Params = gin.Params{{Key: "event_id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.ListByEvent(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNewParticipantHandler(t *testing.T) {
	logger := zap.NewNop()
	handler := NewParticipantHandler(&service.ParticipantService{}, logger)
	assert.NotNil(t, handler)
}

func TestParticipantHandler_Confirm_NoOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/participants/"+uuid.New().String()+"/confirm", nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.Confirm(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_Confirm_InvalidParticipantID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/participants/invalid/confirm", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.Confirm(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_CheckIn_NoOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/participants/"+uuid.New().String()+"/checkin", nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.CheckIn(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_CheckIn_InvalidParticipantID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/participants/invalid/checkin", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.CheckIn(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_BatchCreate_NoOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/"+uuid.New().String()+"/participants/batch", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.BatchCreate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_BatchCreate_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/invalid/participants/batch", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Set("organization_id", uuid.New().String())

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.BatchCreate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParticipantHandler_BatchCreate_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events/"+uuid.New().String()+"/participants/batch", bytes.NewReader([]byte("invalid")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	c.Set("organization_id", uuid.New().String())

	handler := &ParticipantHandler{
		service: nil,
		logger:  logger,
	}

	handler.BatchCreate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
