package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/dto"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"event-coming/internal/service"
)

func setupEventCacheTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, client
}

func TestNewEventCacheHandler(t *testing.T) {
	_, client := setupEventCacheTestRedis(t)
	defer client.Close()

	logger := zap.NewNop()
	svc := service.NewEventCacheService(client)
	handler := NewEventCacheHandler(svc, logger)

	assert.NotNil(t, handler)
}

func TestEventCacheHandler_GetEventCache(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	tests := []struct {
		name           string
		orgID          string
		eventID        string
		setupRedis     func(*miniredis.Miniredis, uuid.UUID, uuid.UUID)
		expectedStatus int
	}{
		{
			name:    "successful get empty cache",
			orgID:   uuid.New().String(),
			eventID: uuid.New().String(),
			setupRedis: func(mr *miniredis.Miniredis, orgID, eventID uuid.UUID) {
				// No data to add
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "successful get with data",
			orgID:   uuid.New().String(),
			eventID: uuid.New().String(),
			setupRedis: func(mr *miniredis.Miniredis, orgID, eventID uuid.UUID) {
				participantID := uuid.New()
				// Add location
				loc := domain.Location{
					ID:            uuid.New(),
					ParticipantID: participantID,
					EventID:       eventID,
					Latitude:      -23.5505,
					Longitude:     -46.6333,
					Timestamp:     time.Now(),
				}
				locJSON, _ := json.Marshal(loc)
				locKey := "location:latest:" + eventID.String() + ":" + participantID.String()
				mr.Set(locKey, string(locJSON))

				// Add confirmation
				conf := dto.ParticipantConfirmationData{
					ParticipantID: participantID,
					Status:        domain.ParticipantStatusConfirmed,
					UpdatedAt:     time.Now(),
				}
				confJSON, _ := json.Marshal(conf)
				confKey := "confirmation:" + orgID.String() + ":" + eventID.String() + ":" + participantID.String()
				mr.Set(confKey, string(confJSON))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid organization id",
			orgID:          "invalid-uuid",
			eventID:        uuid.New().String(),
			setupRedis:     func(mr *miniredis.Miniredis, orgID, eventID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid event id",
			orgID:          uuid.New().String(),
			eventID:        "invalid-uuid",
			setupRedis:     func(mr *miniredis.Miniredis, orgID, eventID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr, client := setupEventCacheTestRedis(t)
			defer mr.Close()
			defer client.Close()

			// Parse UUIDs if valid
			var orgUUID, eventUUID uuid.UUID
			if tt.orgID != "invalid-uuid" {
				orgUUID, _ = uuid.Parse(tt.orgID)
			}
			if tt.eventID != "invalid-uuid" {
				eventUUID, _ = uuid.Parse(tt.eventID)
			}
			tt.setupRedis(mr, orgUUID, eventUUID)

			svc := service.NewEventCacheService(client)
			handler := NewEventCacheHandler(svc, logger)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/"+tt.orgID+"/"+tt.eventID+"/cache", nil)
			c.Params = gin.Params{
				{Key: "organization", Value: tt.orgID},
				{Key: "event", Value: tt.eventID},
			}

			handler.GetEventCache(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestEventCacheHandler_GetLocationsOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	tests := []struct {
		name           string
		orgID          string
		eventID        string
		expectedStatus int
	}{
		{
			name:           "successful get locations",
			orgID:          uuid.New().String(),
			eventID:        uuid.New().String(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid organization id",
			orgID:          "invalid",
			eventID:        uuid.New().String(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid event id",
			orgID:          uuid.New().String(),
			eventID:        "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr, client := setupEventCacheTestRedis(t)
			defer mr.Close()
			defer client.Close()

			svc := service.NewEventCacheService(client)
			handler := NewEventCacheHandler(svc, logger)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/"+tt.orgID+"/"+tt.eventID+"/locations", nil)
			c.Params = gin.Params{
				{Key: "organization", Value: tt.orgID},
				{Key: "event", Value: tt.eventID},
			}

			handler.GetLocationsOnly(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestEventCacheHandler_GetConfirmationsOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	tests := []struct {
		name           string
		orgID          string
		eventID        string
		expectedStatus int
	}{
		{
			name:           "successful get confirmations",
			orgID:          uuid.New().String(),
			eventID:        uuid.New().String(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid organization id",
			orgID:          "invalid",
			eventID:        uuid.New().String(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid event id",
			orgID:          uuid.New().String(),
			eventID:        "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr, client := setupEventCacheTestRedis(t)
			defer mr.Close()
			defer client.Close()

			svc := service.NewEventCacheService(client)
			handler := NewEventCacheHandler(svc, logger)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/"+tt.orgID+"/"+tt.eventID+"/confirmations", nil)
			c.Params = gin.Params{
				{Key: "organization", Value: tt.orgID},
				{Key: "event", Value: tt.eventID},
			}

			handler.GetConfirmationsOnly(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestEventCacheHandler_GetEventCache_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	// Create a broken redis client (pointing to invalid address)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:1", // Invalid port
	})
	defer client.Close()

	// Force context timeout to simulate error
	svc := service.NewEventCacheService(client)
	handler := NewEventCacheHandler(svc, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	orgID := uuid.New().String()
	eventID := uuid.New().String()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/"+orgID+"/"+eventID+"/cache", nil).WithContext(ctx)
	c.Params = gin.Params{
		{Key: "organization", Value: orgID},
		{Key: "event", Value: eventID},
	}

	handler.GetEventCache(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
