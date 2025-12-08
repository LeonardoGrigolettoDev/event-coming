package mocks

import (
	"context"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockLocationBuffer is a mock implementation of LocationBuffer
type MockLocationBuffer struct {
	mock.Mock
}

func (m *MockLocationBuffer) Push(ctx context.Context, location *domain.Location) error {
	args := m.Called(ctx, location)
	return args.Error(0)
}

func (m *MockLocationBuffer) PushWithTTL(ctx context.Context, location *domain.Location, ttl time.Duration) error {
	args := m.Called(ctx, location, ttl)
	return args.Error(0)
}

func (m *MockLocationBuffer) PushWithEventEndTime(ctx context.Context, location *domain.Location, eventEndTime time.Time) error {
	args := m.Called(ctx, location, eventEndTime)
	return args.Error(0)
}

func (m *MockLocationBuffer) GetLatest(ctx context.Context, participantID uuid.UUID) (*domain.Location, error) {
	args := m.Called(ctx, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Location), args.Error(1)
}

func (m *MockLocationBuffer) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]*domain.Location, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Location), args.Error(1)
}

func (m *MockLocationBuffer) GetHistory(ctx context.Context, participantID uuid.UUID, count int) ([]*domain.Location, error) {
	args := m.Called(ctx, participantID, count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Location), args.Error(1)
}

func (m *MockLocationBuffer) Flush(ctx context.Context) ([]*domain.Location, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Location), args.Error(1)
}

func (m *MockLocationBuffer) SetLatestLocation(ctx context.Context, location *domain.Location, eventEndTime time.Time) error {
	args := m.Called(ctx, location, eventEndTime)
	return args.Error(0)
}

// MockRedisCache is a mock implementation of RedisCache
type MockRedisCache struct {
	mock.Mock
}

func (m *MockRedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockRedisCache) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisCache) LPush(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockRedisCache) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	args := m.Called(ctx, key, start, stop)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisCache) LTrim(ctx context.Context, key string, start, stop int64) error {
	args := m.Called(ctx, key, start, stop)
	return args.Error(0)
}

func (m *MockRedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	args := m.Called(ctx, key, expiration)
	return args.Error(0)
}

func (m *MockRedisCache) SAdd(ctx context.Context, key string, members ...interface{}) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisCache) SMembers(ctx context.Context, key string) ([]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisCache) SRem(ctx context.Context, key string, members ...interface{}) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisCache) Publish(ctx context.Context, channel string, message interface{}) error {
	args := m.Called(ctx, channel, message)
	return args.Error(0)
}

// MockEventCacheService is a mock implementation of EventCacheService
type MockEventCacheService struct {
	mock.Mock
}

func (m *MockEventCacheService) CacheEvent(ctx context.Context, event *domain.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventCacheService) GetCachedEvent(ctx context.Context, eventID uuid.UUID) (*dto.EventCacheResponse, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EventCacheResponse), args.Error(1)
}

func (m *MockEventCacheService) InvalidateEvent(ctx context.Context, eventID uuid.UUID) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

func (m *MockEventCacheService) CacheParticipantLocations(ctx context.Context, eventID uuid.UUID, locations []*dto.LocationResponse) error {
	args := m.Called(ctx, eventID, locations)
	return args.Error(0)
}

func (m *MockEventCacheService) GetCachedParticipantLocations(ctx context.Context, eventID uuid.UUID) ([]*dto.LocationResponse, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.LocationResponse), args.Error(1)
}

// MockPubSub is a mock implementation of PubSub
type MockPubSub struct {
	mock.Mock
}

func (m *MockPubSub) Publish(ctx context.Context, channel string, message interface{}) error {
	args := m.Called(ctx, channel, message)
	return args.Error(0)
}

func (m *MockPubSub) Subscribe(ctx context.Context, channel string) (<-chan string, error) {
	args := m.Called(ctx, channel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(<-chan string), args.Error(1)
}

func (m *MockPubSub) Unsubscribe(ctx context.Context, channel string) error {
	args := m.Called(ctx, channel)
	return args.Error(0)
}
