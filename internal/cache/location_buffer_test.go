package cache_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"event-coming/internal/cache"
	"event-coming/internal/domain"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	s, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	return s, client
}

func TestNewLocationBuffer(t *testing.T) {
	_, client := setupTestRedis(t)
	defer client.Close()

	buffer := cache.NewLocationBuffer(client)
	assert.NotNil(t, buffer)
}

func TestLocationBuffer_Push(t *testing.T) {
	s, client := setupTestRedis(t)
	defer s.Close()
	defer client.Close()

	buffer := cache.NewLocationBuffer(client)
	ctx := context.Background()

	location := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: uuid.New(),
		EventID:       uuid.New(),
		EntityID:      uuid.New(),
		Latitude:      -23.550520,
		Longitude:     -46.633308,
		Timestamp:     time.Now(),
		CreatedAt:     time.Now(),
	}

	err := buffer.Push(ctx, location)
	assert.NoError(t, err)

	// Verify location is in buffer
	bufferKey := "location:buffer:" + location.EntityID.String()
	assert.True(t, s.Exists(bufferKey))
}

func TestLocationBuffer_PushWithEventEndTime(t *testing.T) {
	s, client := setupTestRedis(t)
	defer s.Close()
	defer client.Close()

	buffer := cache.NewLocationBuffer(client)
	ctx := context.Background()

	location := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: uuid.New(),
		EventID:       uuid.New(),
		EntityID:      uuid.New(),
		Latitude:      -23.550520,
		Longitude:     -46.633308,
		Timestamp:     time.Now(),
		CreatedAt:     time.Now(),
	}

	// Event ends in 2 hours
	eventEndTime := time.Now().Add(2 * time.Hour)

	err := buffer.PushWithEventEndTime(ctx, location, eventEndTime)
	assert.NoError(t, err)

	// Verify location is in cache
	cacheKey := "location:latest:" + location.EventID.String() + ":" + location.ParticipantID.String()
	assert.True(t, s.Exists(cacheKey))
}

func TestLocationBuffer_PushWithEventEndTime_EventEnded(t *testing.T) {
	s, client := setupTestRedis(t)
	defer s.Close()
	defer client.Close()

	buffer := cache.NewLocationBuffer(client)
	ctx := context.Background()

	location := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: uuid.New(),
		EventID:       uuid.New(),
		EntityID:      uuid.New(),
		Latitude:      -23.550520,
		Longitude:     -46.633308,
		Timestamp:     time.Now(),
		CreatedAt:     time.Now(),
	}

	// Event ended 1 hour ago
	eventEndTime := time.Now().Add(-1 * time.Hour)

	err := buffer.PushWithEventEndTime(ctx, location, eventEndTime)
	assert.NoError(t, err)

	// Verify location is still cached (with minimum TTL)
	cacheKey := "location:latest:" + location.EventID.String() + ":" + location.ParticipantID.String()
	assert.True(t, s.Exists(cacheKey))
}

func TestLocationBuffer_PushWithTTL(t *testing.T) {
	s, client := setupTestRedis(t)
	defer s.Close()
	defer client.Close()

	buffer := cache.NewLocationBuffer(client)
	ctx := context.Background()

	location := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: uuid.New(),
		EventID:       uuid.New(),
		EntityID:      uuid.New(),
		Latitude:      -23.550520,
		Longitude:     -46.633308,
		Timestamp:     time.Now(),
		CreatedAt:     time.Now(),
	}

	ttl := 1 * time.Hour
	err := buffer.PushWithTTL(ctx, location, ttl)
	assert.NoError(t, err)
}

func TestLocationBuffer_SetLatestLocation(t *testing.T) {
	s, client := setupTestRedis(t)
	defer s.Close()
	defer client.Close()

	buffer := cache.NewLocationBuffer(client)
	ctx := context.Background()

	location := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: uuid.New(),
		EventID:       uuid.New(),
		EntityID:      uuid.New(),
		Latitude:      -23.550520,
		Longitude:     -46.633308,
		Timestamp:     time.Now(),
		CreatedAt:     time.Now(),
	}

	eventEndTime := time.Now().Add(3 * time.Hour)

	err := buffer.SetLatestLocation(ctx, location, eventEndTime)
	assert.NoError(t, err)

	// Verify location is cached
	cacheKey := "location:latest:" + location.EventID.String() + ":" + location.ParticipantID.String()
	assert.True(t, s.Exists(cacheKey))

	// Verify cached data is correct
	data, err := client.Get(ctx, cacheKey).Result()
	assert.NoError(t, err)

	var cached domain.Location
	err = json.Unmarshal([]byte(data), &cached)
	assert.NoError(t, err)
	assert.Equal(t, location.ID, cached.ID)
	assert.Equal(t, location.Latitude, cached.Latitude)
	assert.Equal(t, location.Longitude, cached.Longitude)
}

func TestLocationBuffer_SetLatestLocation_Update(t *testing.T) {
	s, client := setupTestRedis(t)
	defer s.Close()
	defer client.Close()

	buffer := cache.NewLocationBuffer(client)
	ctx := context.Background()

	eventID := uuid.New()
	participantID := uuid.New()
	entityID := uuid.New()
	eventEndTime := time.Now().Add(3 * time.Hour)

	// First location
	location1 := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: participantID,
		EventID:       eventID,
		EntityID:      entityID,
		Latitude:      -23.550520,
		Longitude:     -46.633308,
		Timestamp:     time.Now(),
		CreatedAt:     time.Now(),
	}

	err := buffer.SetLatestLocation(ctx, location1, eventEndTime)
	assert.NoError(t, err)

	// Update with new location
	location2 := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: participantID,
		EventID:       eventID,
		EntityID:      entityID,
		Latitude:      -23.560520, // New position
		Longitude:     -46.643308,
		Timestamp:     time.Now().Add(1 * time.Minute),
		CreatedAt:     time.Now(),
	}

	err = buffer.SetLatestLocation(ctx, location2, eventEndTime)
	assert.NoError(t, err)

	// Verify latest location is location2
	retrieved, err := buffer.GetLatestLocation(ctx, eventID, participantID)
	assert.NoError(t, err)
	assert.Equal(t, location2.ID, retrieved.ID)
	assert.Equal(t, location2.Latitude, retrieved.Latitude)
}

func TestLocationBuffer_GetLatestLocation(t *testing.T) {
	s, client := setupTestRedis(t)
	defer s.Close()
	defer client.Close()

	buffer := cache.NewLocationBuffer(client)
	ctx := context.Background()

	eventID := uuid.New()
	participantID := uuid.New()
	entityID := uuid.New()

	location := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: participantID,
		EventID:       eventID,
		EntityID:      entityID,
		Latitude:      -23.550520,
		Longitude:     -46.633308,
		Timestamp:     time.Now(),
		CreatedAt:     time.Now(),
	}

	// Store location
	err := buffer.SetLatestLocation(ctx, location, time.Now().Add(3*time.Hour))
	assert.NoError(t, err)

	// Retrieve location
	retrieved, err := buffer.GetLatestLocation(ctx, eventID, participantID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, location.ID, retrieved.ID)
	assert.Equal(t, location.Latitude, retrieved.Latitude)
	assert.Equal(t, location.Longitude, retrieved.Longitude)
}

func TestLocationBuffer_GetLatestLocation_NotFound(t *testing.T) {
	s, client := setupTestRedis(t)
	defer s.Close()
	defer client.Close()

	buffer := cache.NewLocationBuffer(client)
	ctx := context.Background()

	// Try to get non-existent location
	retrieved, err := buffer.GetLatestLocation(ctx, uuid.New(), uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestLocationBuffer_GetLatestLocationsForEvent(t *testing.T) {
	s, client := setupTestRedis(t)
	defer s.Close()
	defer client.Close()

	buffer := cache.NewLocationBuffer(client)
	ctx := context.Background()

	eventID := uuid.New()
	entityID := uuid.New()
	participant1 := uuid.New()
	participant2 := uuid.New()
	eventEndTime := time.Now().Add(3 * time.Hour)

	// Store locations for 2 participants
	loc1 := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: participant1,
		EventID:       eventID,
		EntityID:      entityID,
		Latitude:      -23.550520,
		Longitude:     -46.633308,
		Timestamp:     time.Now(),
		CreatedAt:     time.Now(),
	}

	loc2 := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: participant2,
		EventID:       eventID,
		EntityID:      entityID,
		Latitude:      -22.906847,
		Longitude:     -43.172897,
		Timestamp:     time.Now(),
		CreatedAt:     time.Now(),
	}

	err := buffer.SetLatestLocation(ctx, loc1, eventEndTime)
	assert.NoError(t, err)
	err = buffer.SetLatestLocation(ctx, loc2, eventEndTime)
	assert.NoError(t, err)

	// Retrieve all locations
	participantIDs := []uuid.UUID{participant1, participant2}
	locations, err := buffer.GetLatestLocationsForEvent(ctx, eventID, participantIDs)
	assert.NoError(t, err)
	assert.Len(t, locations, 2)
}

func TestLocationBuffer_GetLatestLocationsForEvent_Empty(t *testing.T) {
	s, client := setupTestRedis(t)
	defer s.Close()
	defer client.Close()

	buffer := cache.NewLocationBuffer(client)
	ctx := context.Background()

	// Empty participant IDs
	locations, err := buffer.GetLatestLocationsForEvent(ctx, uuid.New(), []uuid.UUID{})
	assert.NoError(t, err)
	assert.Empty(t, locations)
}

func TestLocationBuffer_PopBatch(t *testing.T) {
	s, client := setupTestRedis(t)
	defer s.Close()
	defer client.Close()

	buffer := cache.NewLocationBuffer(client)
	ctx := context.Background()

	entityID := uuid.New()
	eventID := uuid.New()

	// Push multiple locations
	for i := 0; i < 5; i++ {
		location := &domain.Location{
			ID:            uuid.New(),
			ParticipantID: uuid.New(),
			EventID:       eventID,
			EntityID:      entityID,
			Latitude:      -23.550520 + float64(i)*0.001,
			Longitude:     -46.633308,
			Timestamp:     time.Now(),
			CreatedAt:     time.Now(),
		}
		err := buffer.Push(ctx, location)
		assert.NoError(t, err)
	}

	// Pop batch of 3
	locations, err := buffer.PopBatch(ctx, entityID, 3)
	assert.NoError(t, err)
	assert.Len(t, locations, 3)

	// Pop remaining 2
	locations, err = buffer.PopBatch(ctx, entityID, 10)
	assert.NoError(t, err)
	assert.Len(t, locations, 2)

	// Pop empty buffer
	locations, err = buffer.PopBatch(ctx, entityID, 10)
	assert.NoError(t, err)
	assert.Empty(t, locations)
}
