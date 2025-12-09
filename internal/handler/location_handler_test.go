package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"event-coming/internal/service"
	"event-coming/internal/service/eta"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLocationHandler_CreateLocation_InvalidParticipantID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/participants/invalid/locations", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.CreateLocation(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLocationHandler_CreateLocation_NoEntityID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/participants/"+uuid.New().String()+"/locations", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	// No entity_id in context

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.CreateLocation(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLocationHandler_CreateLocation_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/participants/"+uuid.New().String()+"/locations", bytes.NewReader([]byte("invalid")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	c.Set("entity_id", uuid.New())

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.CreateLocation(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLocationHandler_GetLatestLocation_InvalidParticipantID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/participants/invalid/location", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.GetLatestLocation(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLocationHandler_GetLocationHistory_InvalidParticipantID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/participants/invalid/locations", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.GetLocationHistory(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLocationHandler_GetLocationsByEvent_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/events/invalid/locations", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.GetEventLocations(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLocationHandler_GetParticipantETA_InvalidParticipantID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/participants/invalid/eta", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.GetParticipantETA(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLocationHandler_GetEventETAs_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/events/invalid/etas", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.GetEventETAs(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNewLocationHandler(t *testing.T) {
	handler := NewLocationHandler(&service.LocationService{}, &eta.ETAService{}, &service.EventService{})
	assert.NotNil(t, handler)
}

func TestLocationHandler_GetLatestLocation_NoEntityID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/participants/"+uuid.New().String()+"/location", nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	// No entity_id in context

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.GetLatestLocation(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLocationHandler_GetLocationHistory_NoEntityID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/participants/"+uuid.New().String()+"/locations", nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	// No entity_id in context

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.GetLocationHistory(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLocationHandler_GetEventLocations_NoEntityID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/events/"+uuid.New().String()+"/locations", nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	// No entity_id in context

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.GetEventLocations(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLocationHandler_GetParticipantETA_NoEntityID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/participants/"+uuid.New().String()+"/eta", nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	// No entity_id in context

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.GetParticipantETA(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLocationHandler_GetEventETAs_NoEntityID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/events/"+uuid.New().String()+"/etas", nil)
	c.Params = gin.Params{{Key: "id", Value: uuid.New().String()}}
	// No entity_id in context

	handler := &LocationHandler{
		locationService: nil,
		etaService:      nil,
		eventService:    nil,
	}

	handler.GetEventETAs(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
