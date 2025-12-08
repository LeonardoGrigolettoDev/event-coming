package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// LocationBuffer handles buffering of location data in Redis
type LocationBuffer struct {
	client *redis.Client
}

// NewLocationBuffer creates a new location buffer
func NewLocationBuffer(client *redis.Client) *LocationBuffer {
	return &LocationBuffer{client: client}
}

// Push adds a location to the buffer (uses default 24h TTL)
func (b *LocationBuffer) Push(ctx context.Context, location *domain.Location) error {
	return b.PushWithTTL(ctx, location, 24*time.Hour)
}

// PushWithEventEndTime adds a location to the buffer with TTL based on event end time
func (b *LocationBuffer) PushWithEventEndTime(ctx context.Context, location *domain.Location, eventEndTime time.Time) error {
	// Calculate TTL based on event end time
	ttl := time.Until(eventEndTime)
	if ttl <= 0 {
		// Event already ended, use minimum TTL of 1 hour for cleanup
		ttl = 1 * time.Hour
	}
	// Add 1 hour buffer after event ends
	ttl += 1 * time.Hour

	return b.PushWithTTL(ctx, location, ttl)
}

// PushWithTTL adds a location to the buffer with custom TTL
func (b *LocationBuffer) PushWithTTL(ctx context.Context, location *domain.Location, ttl time.Duration) error {
	// Serialize location
	data, err := json.Marshal(location)
	if err != nil {
		return fmt.Errorf("failed to marshal location: %w", err)
	}

	// Add to list buffer
	bufferKey := fmt.Sprintf("location:buffer:%s", location.EntityID)
	if err := b.client.RPush(ctx, bufferKey, data).Err(); err != nil {
		return fmt.Errorf("failed to push to buffer: %w", err)
	}

	// Update latest location cache with TTL
	cacheKey := fmt.Sprintf("location:latest:%s:%s", location.EventID, location.ParticipantID)
	if err := b.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to cache latest location: %w", err)
	}

	// Publish to pub/sub for real-time updates
	channel := fmt.Sprintf("location:updates:%s", location.EventID)
	if err := b.client.Publish(ctx, channel, data).Err(); err != nil {
		// Log error but don't fail
		fmt.Printf("failed to publish location update: %v\n", err)
	}

	return nil
}

// SetLatestLocation sets or updates the latest location with TTL based on event end time
func (b *LocationBuffer) SetLatestLocation(ctx context.Context, location *domain.Location, eventEndTime time.Time) error {
	// Serialize location
	data, err := json.Marshal(location)
	if err != nil {
		return fmt.Errorf("failed to marshal location: %w", err)
	}

	// Calculate TTL based on event end time
	ttl := time.Until(eventEndTime)
	if ttl <= 0 {
		ttl = 1 * time.Hour
	}
	ttl += 1 * time.Hour // Add buffer after event ends

	cacheKey := fmt.Sprintf("location:latest:%s:%s", location.EventID, location.ParticipantID)

	// Use SET with TTL - this creates if not exists or updates if exists
	if err := b.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set latest location: %w", err)
	}

	// Also publish for real-time updates
	channel := fmt.Sprintf("location:updates:%s", location.EventID)
	b.client.Publish(ctx, channel, data)

	return nil
}

// PopBatch retrieves and removes a batch of locations from the buffer
func (b *LocationBuffer) PopBatch(ctx context.Context, orgID uuid.UUID, batchSize int) ([]*domain.Location, error) {
	bufferKey := fmt.Sprintf("location:buffer:%s", orgID)

	// Use Lua script for atomic pop
	script := redis.NewScript(`
		local key = KEYS[1]
		local count = tonumber(ARGV[1])
		local items = redis.call('LRANGE', key, 0, count - 1)
		if #items > 0 then
			redis.call('LTRIM', key, #items, -1)
		end
		return items
	`)

	result, err := script.Run(ctx, b.client, []string{bufferKey}, batchSize).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to pop batch: %w", err)
	}

	items, ok := result.([]interface{})
	if !ok {
		return []*domain.Location{}, nil
	}

	var locations []*domain.Location
	for _, item := range items {
		str, ok := item.(string)
		if !ok {
			continue
		}

		var loc domain.Location
		if err := json.Unmarshal([]byte(str), &loc); err != nil {
			continue
		}
		locations = append(locations, &loc)
	}

	return locations, nil
}

// GetLatestLocation retrieves the latest location for a participant
func (b *LocationBuffer) GetLatestLocation(ctx context.Context, eventID, participantID uuid.UUID) (*domain.Location, error) {
	cacheKey := fmt.Sprintf("location:latest:%s:%s", eventID, participantID)

	data, err := b.client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest location: %w", err)
	}

	var location domain.Location
	if err := json.Unmarshal([]byte(data), &location); err != nil {
		return nil, fmt.Errorf("failed to unmarshal location: %w", err)
	}

	return &location, nil
}

// GetLatestLocationsForEvent retrieves all latest locations for an event
func (b *LocationBuffer) GetLatestLocationsForEvent(ctx context.Context, eventID uuid.UUID, participantIDs []uuid.UUID) ([]*domain.Location, error) {
	if len(participantIDs) == 0 {
		return []*domain.Location{}, nil
	}

	// Build keys
	keys := make([]string, len(participantIDs))
	for i, pid := range participantIDs {
		keys[i] = fmt.Sprintf("location:latest:%s:%s", eventID, pid)
	}

	// Use MGET for batch retrieval
	results, err := b.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	var locations []*domain.Location
	for _, result := range results {
		if result == nil {
			continue
		}

		str, ok := result.(string)
		if !ok {
			continue
		}

		var loc domain.Location
		if err := json.Unmarshal([]byte(str), &loc); err != nil {
			continue
		}
		locations = append(locations, &loc)
	}

	return locations, nil
}

// SubscribeToEvent subscribes to location updates for an event
func (b *LocationBuffer) SubscribeToEvent(ctx context.Context, eventID uuid.UUID) *redis.PubSub {
	channel := fmt.Sprintf("location:updates:%s", eventID)
	return b.client.Subscribe(ctx, channel)
}
