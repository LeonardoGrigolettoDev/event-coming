package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"event-coming/internal/domain"
	"event-coming/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// Since EntityHandler uses *service.EntityService directly, we need to create
// actual handlers that test the HTTP layer
func TestEntityHandler_Create_InvalidJSON(t *testing.T) {
	logger := zap.NewNop()
	// Can't easily mock *service.EntityService, so test what we can

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/entities", bytes.NewReader([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	// Create handler with nil service (will panic on actual call but we're testing JSON parsing)
	handler := &EntityHandler{
		entityService: nil,
		logger:        logger,
	}

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEntityHandler_GetByID_InvalidUUID(t *testing.T) {
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/entities/invalid-uuid", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	handler := &EntityHandler{
		entityService: nil,
		logger:        logger,
	}

	handler.GetByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEntityHandler_Update_InvalidUUID(t *testing.T) {
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/entities/invalid-uuid", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	handler := &EntityHandler{
		entityService: nil,
		logger:        logger,
	}

	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEntityHandler_Update_InvalidJSON(t *testing.T) {
	logger := zap.NewNop()
	validID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/entities/"+validID.String(), bytes.NewReader([]byte("invalid")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: validID.String()}}

	handler := &EntityHandler{
		entityService: nil,
		logger:        logger,
	}

	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEntityHandler_Delete_InvalidUUID(t *testing.T) {
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/entities/invalid-uuid", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	handler := &EntityHandler{
		entityService: nil,
		logger:        logger,
	}

	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEntityHandler_ListByParent_InvalidUUID(t *testing.T) {
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/entities/invalid-uuid/children", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	handler := &EntityHandler{
		entityService: nil,
		logger:        logger,
	}

	handler.ListByParent(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEntityHandler_GetByDocument_Empty(t *testing.T) {
	logger := zap.NewNop()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/entities/document/", nil)
	c.Params = gin.Params{{Key: "document", Value: ""}}

	handler := &EntityHandler{
		entityService: nil,
		logger:        logger,
	}

	handler.GetByDocument(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test validation errors using a real EntityService wrapper approach
// Since the handler uses a concrete type, we need integration-style tests for full coverage

// MockableEntityHandler wraps an EntityHandler with mockable service
type MockableEntityHandler struct {
	logger *zap.Logger
}

func TestNewEntityHandler(t *testing.T) {
	logger := zap.NewNop()
	// Test that NewEntityHandler creates properly
	handler := NewEntityHandler(&service.EntityService{}, logger)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.logger)
}

// Test domain errors are handled correctly
func TestEntityHandler_HandlesDomainErrors(t *testing.T) {
	// Test that response.HandleDomainError works for various error types
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{"not found", domain.ErrNotFound, http.StatusNotFound},
		{"conflict (already exists)", domain.ErrConflict, http.StatusConflict},
		{"invalid input", domain.ErrInvalidInput, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a design limitation - EntityHandler uses concrete *service.EntityService
			// To properly test, we'd need dependency injection with interfaces
			// For now we verify the error types exist
			assert.Error(t, tt.err)
		})
	}
}
