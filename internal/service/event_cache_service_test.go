package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/dto"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupEventCacheTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, client
}

func TestNewEventCacheService(t *testing.T) {
	_, client := setupEventCacheTestRedis(t)
	defer client.Close()

	svc := NewEventCacheService(client)
	assert.NotNil(t, svc)
}

func TestEventCacheService_GetEventCacheData(t *testing.T) {
	ctx := context.Background()
	entID := uuid.New()
	eventID := uuid.New()
	participantID := uuid.New()

	t.Run("empty cache returns empty data", func(t *testing.T) {
		mr, client := setupEventCacheTestRedis(t)
		defer mr.Close()
		defer client.Close()

		svc := NewEventCacheService(client)
		data, err := svc.GetEventCacheData(ctx, entID, eventID)

		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, entID, data.EntityID)
		assert.Equal(t, eventID, data.EventID)
		assert.Empty(t, data.Locations)
		assert.Empty(t, data.Confirmations)
		assert.Equal(t, 0, data.TotalLocations)
	})

	t.Run("with locations and confirmations", func(t *testing.T) {
		mr, client := setupEventCacheTestRedis(t)
		defer mr.Close()
		defer client.Close()

		// Add location data
		loc := domain.Location{
			ID:            uuid.New(),
			ParticipantID: participantID,
			EventID:       eventID,
			Latitude:      -23.5505,
			Longitude:     -46.6333,
			Accuracy:      floatPtr(10.0),
			Speed:         floatPtr(5.0),
			Heading:       floatPtr(180.0),
			Timestamp:     time.Now(),
		}
		locJSON, _ := json.Marshal(loc)
		locKey := "location:latest:" + eventID.String() + ":" + participantID.String()
		mr.Set(locKey, string(locJSON))

		// Add confirmation data
		conf := dto.ParticipantConfirmationData{
			ParticipantID: participantID,
			Status:        domain.ParticipantStatusConfirmed,
			ConfirmedAt:   timePtr(time.Now()),
			UpdatedAt:     time.Now(),
		}
		confJSON, _ := json.Marshal(conf)
		confKey := "confirmation:" + entID.String() + ":" + eventID.String() + ":" + participantID.String()
		mr.Set(confKey, string(confJSON))

		svc := NewEventCacheService(client)
		data, err := svc.GetEventCacheData(ctx, entID, eventID)

		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Len(t, data.Locations, 1)
		assert.Len(t, data.Confirmations, 1)
		assert.Equal(t, 1, data.TotalConfirmed)
		assert.Equal(t, 0, data.TotalPending)
		assert.Equal(t, 0, data.TotalDenied)
	})

	t.Run("with various confirmation statuses", func(t *testing.T) {
		mr, client := setupEventCacheTestRedis(t)
		defer mr.Close()
		defer client.Close()

		// Add pending confirmation
		pending := dto.ParticipantConfirmationData{
			ParticipantID: uuid.New(),
			Status:        domain.ParticipantStatusPending,
			UpdatedAt:     time.Now(),
		}
		pendingJSON, _ := json.Marshal(pending)
		pendingKey := "confirmation:" + entID.String() + ":" + eventID.String() + ":" + pending.ParticipantID.String()
		mr.Set(pendingKey, string(pendingJSON))

		// Add denied confirmation
		denied := dto.ParticipantConfirmationData{
			ParticipantID: uuid.New(),
			Status:        domain.ParticipantStatusDenied,
			UpdatedAt:     time.Now(),
		}
		deniedJSON, _ := json.Marshal(denied)
		deniedKey := "confirmation:" + entID.String() + ":" + eventID.String() + ":" + denied.ParticipantID.String()
		mr.Set(deniedKey, string(deniedJSON))

		// Add checked-in confirmation
		checkedIn := dto.ParticipantConfirmationData{
			ParticipantID: uuid.New(),
			Status:        domain.ParticipantStatusCheckedIn,
			CheckedInAt:   timePtr(time.Now()),
			UpdatedAt:     time.Now(),
		}
		checkedInJSON, _ := json.Marshal(checkedIn)
		checkedInKey := "confirmation:" + entID.String() + ":" + eventID.String() + ":" + checkedIn.ParticipantID.String()
		mr.Set(checkedInKey, string(checkedInJSON))

		svc := NewEventCacheService(client)
		data, err := svc.GetEventCacheData(ctx, entID, eventID)

		assert.NoError(t, err)
		assert.Len(t, data.Confirmations, 3)
		assert.Equal(t, 1, data.TotalConfirmed) // checked-in counts as confirmed
		assert.Equal(t, 1, data.TotalPending)
		assert.Equal(t, 1, data.TotalDenied)
	})
}

func TestEventCacheService_SetConfirmation(t *testing.T) {
	ctx := context.Background()
	entID := uuid.New()
	eventID := uuid.New()
	now := time.Now()

	tests := []struct {
		name        string
		participant *domain.Participant
		wantErr     bool
	}{
		{
			name: "set confirmed participant",
			participant: &domain.Participant{
				ID:          uuid.New(),
				EntityID:    entID,
				EventID:     eventID,
				Status:      domain.ParticipantStatusConfirmed,
				ConfirmedAt: &now,
			},
			wantErr: false,
		},
		{
			name: "set checked-in participant",
			participant: &domain.Participant{
				ID:          uuid.New(),
				EntityID:    entID,
				EventID:     eventID,
				Status:      domain.ParticipantStatusCheckedIn,
				ConfirmedAt: &now,
				CheckedInAt: &now,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr, client := setupEventCacheTestRedis(t)
			defer mr.Close()
			defer client.Close()

			svc := NewEventCacheService(client)
			err := svc.SetConfirmation(ctx, entID, eventID, tt.participant)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify stored data
				key := "confirmation:" + entID.String() + ":" + eventID.String() + ":" + tt.participant.ID.String()
				val, getErr := mr.Get(key)
				assert.NoError(t, getErr)
				assert.NotEmpty(t, val)

				var stored dto.ParticipantConfirmationData
				err := json.Unmarshal([]byte(val), &stored)
				assert.NoError(t, err)
				assert.Equal(t, tt.participant.ID, stored.ParticipantID)
				assert.Equal(t, tt.participant.Status, stored.Status)
			}
		})
	}
}

func TestEventCacheService_DeleteConfirmation(t *testing.T) {
	ctx := context.Background()
	entID := uuid.New()
	eventID := uuid.New()
	participantID := uuid.New()

	t.Run("delete existing confirmation", func(t *testing.T) {
		mr, client := setupEventCacheTestRedis(t)
		defer mr.Close()
		defer client.Close()

		// First set a confirmation
		key := "confirmation:" + entID.String() + ":" + eventID.String() + ":" + participantID.String()
		mr.Set(key, `{"participant_id":"test"}`)

		svc := NewEventCacheService(client)
		err := svc.DeleteConfirmation(ctx, entID, eventID, participantID)

		assert.NoError(t, err)
		assert.False(t, mr.Exists(key))
	})

	t.Run("delete non-existing confirmation - no error", func(t *testing.T) {
		mr, client := setupEventCacheTestRedis(t)
		defer mr.Close()
		defer client.Close()

		svc := NewEventCacheService(client)
		err := svc.DeleteConfirmation(ctx, entID, eventID, uuid.New())

		assert.NoError(t, err)
	})
}

func TestEventCacheService_GetLocationsSummary(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()

	t.Run("no locations returns zero", func(t *testing.T) {
		mr, client := setupEventCacheTestRedis(t)
		defer mr.Close()
		defer client.Close()

		svc := NewEventCacheService(client)
		count, err := svc.GetLocationsSummary(ctx, eventID)

		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("counts all locations", func(t *testing.T) {
		mr, client := setupEventCacheTestRedis(t)
		defer mr.Close()
		defer client.Close()

		// Add multiple locations
		for i := 0; i < 5; i++ {
			key := "location:latest:" + eventID.String() + ":" + uuid.New().String()
			mr.Set(key, `{"latitude":-23.5505,"longitude":-46.6333}`)
		}

		svc := NewEventCacheService(client)
		count, err := svc.GetLocationsSummary(ctx, eventID)

		assert.NoError(t, err)
		assert.Equal(t, 5, count)
	})

	t.Run("only counts matching event locations", func(t *testing.T) {
		mr, client := setupEventCacheTestRedis(t)
		defer mr.Close()
		defer client.Close()

		// Add locations for our event
		for i := 0; i < 3; i++ {
			key := "location:latest:" + eventID.String() + ":" + uuid.New().String()
			mr.Set(key, `{"latitude":-23.5505}`)
		}

		// Add locations for another event
		otherEventID := uuid.New()
		for i := 0; i < 2; i++ {
			key := "location:latest:" + otherEventID.String() + ":" + uuid.New().String()
			mr.Set(key, `{"latitude":-23.5505}`)
		}

		svc := NewEventCacheService(client)
		count, err := svc.GetLocationsSummary(ctx, eventID)

		assert.NoError(t, err)
		assert.Equal(t, 3, count)
	})
}

func TestEventCacheService_GetLocations_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	entID := uuid.New()
	eventID := uuid.New()

	mr, client := setupEventCacheTestRedis(t)
	defer mr.Close()
	defer client.Close()

	// Add invalid JSON
	key := "location:latest:" + eventID.String() + ":" + uuid.New().String()
	mr.Set(key, "invalid json")

	// Add valid JSON
	loc := domain.Location{
		ID:            uuid.New(),
		ParticipantID: uuid.New(),
		EventID:       eventID,
		Latitude:      -23.5505,
		Longitude:     -46.6333,
		Timestamp:     time.Now(),
	}
	locJSON, _ := json.Marshal(loc)
	validKey := "location:latest:" + eventID.String() + ":" + loc.ParticipantID.String()
	mr.Set(validKey, string(locJSON))

	svc := NewEventCacheService(client)
	data, err := svc.GetEventCacheData(ctx, entID, eventID)

	// Should not error, just skip invalid entries
	assert.NoError(t, err)
	assert.Len(t, data.Locations, 1)
}

func TestEventCacheService_GetConfirmations_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	entID := uuid.New()
	eventID := uuid.New()

	mr, client := setupEventCacheTestRedis(t)
	defer mr.Close()
	defer client.Close()

	// Add invalid JSON
	key := "confirmation:" + entID.String() + ":" + eventID.String() + ":" + uuid.New().String()
	mr.Set(key, "invalid json")

	// Add valid confirmation
	conf := dto.ParticipantConfirmationData{
		ParticipantID: uuid.New(),
		Status:        domain.ParticipantStatusConfirmed,
		UpdatedAt:     time.Now(),
	}
	confJSON, _ := json.Marshal(conf)
	validKey := "confirmation:" + entID.String() + ":" + eventID.String() + ":" + conf.ParticipantID.String()
	mr.Set(validKey, string(confJSON))

	svc := NewEventCacheService(client)
	data, err := svc.GetEventCacheData(ctx, entID, eventID)

	// Should not error, just skip invalid entries
	assert.NoError(t, err)
	assert.Len(t, data.Confirmations, 1)
}

// Helper functions
func floatPtr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}
