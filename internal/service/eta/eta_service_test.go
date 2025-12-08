package eta_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"event-coming/internal/config"
	"event-coming/internal/domain"
	"event-coming/internal/service/eta"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLocationRepository is a mock implementation of repository.LocationRepository
type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) Create(ctx context.Context, location *domain.Location) error {
	args := m.Called(ctx, location)
	return args.Error(0)
}

func (m *MockLocationRepository) BatchCreate(ctx context.Context, locations []*domain.Location) error {
	args := m.Called(ctx, locations)
	return args.Error(0)
}

func (m *MockLocationRepository) GetLatestByParticipant(ctx context.Context, participantID, entityID uuid.UUID) (*domain.Location, error) {
	args := m.Called(ctx, participantID, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Location), args.Error(1)
}

func (m *MockLocationRepository) GetLatestByEvent(ctx context.Context, eventID, entityID uuid.UUID) ([]*domain.Location, error) {
	args := m.Called(ctx, eventID, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Location), args.Error(1)
}

func (m *MockLocationRepository) GetHistory(ctx context.Context, participantID, entityID uuid.UUID, from, to time.Time) ([]*domain.Location, error) {
	args := m.Called(ctx, participantID, entityID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Location), args.Error(1)
}

func TestNewETAService(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	cfg := &config.OSRMConfig{Enabled: false}

	service := eta.NewETAService(mockRepo, cfg)
	assert.NotNil(t, service)
}

func TestETAService_CalculateETA_Success(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	cfg := &config.OSRMConfig{Enabled: false}
	service := eta.NewETAService(mockRepo, cfg)
	ctx := context.Background()

	participantID := uuid.New()
	entityID := uuid.New()
	targetLat := -23.550520
	targetLng := -46.633308
	now := time.Now()

	// Latest location ~5 km away
	latestLocation := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: participantID,
		EntityID:      entityID,
		Latitude:      -23.500520, // ~5km north
		Longitude:     -46.633308,
		Timestamp:     now,
	}

	// History with 2 locations for velocity calculation
	history := []*domain.Location{
		{
			ID:            uuid.New(),
			ParticipantID: participantID,
			EntityID:      entityID,
			Latitude:      -23.510520,
			Longitude:     -46.633308,
			Timestamp:     now.Add(-5 * time.Minute),
		},
		{
			ID:            uuid.New(),
			ParticipantID: participantID,
			EntityID:      entityID,
			Latitude:      -23.500520,
			Longitude:     -46.633308,
			Timestamp:     now,
		},
	}

	mockRepo.On("GetLatestByParticipant", ctx, participantID, entityID).Return(latestLocation, nil)
	mockRepo.On("GetHistory", ctx, participantID, entityID, mock.Anything, mock.Anything).Return(history, nil)

	result, err := service.CalculateETA(ctx, participantID, entityID, targetLat, targetLng)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, participantID, result.ParticipantID)
	assert.Greater(t, result.DistanceMeters, 0.0)
	assert.Equal(t, "velocity", result.Method)
	mockRepo.AssertExpectations(t)
}

func TestETAService_CalculateETA_SimpleMethod(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	cfg := &config.OSRMConfig{Enabled: false}
	service := eta.NewETAService(mockRepo, cfg)
	ctx := context.Background()

	participantID := uuid.New()
	entityID := uuid.New()
	targetLat := -23.550520
	targetLng := -46.633308
	now := time.Now()

	latestLocation := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: participantID,
		EntityID:      entityID,
		Latitude:      -23.500520,
		Longitude:     -46.633308,
		Timestamp:     now,
	}

	// Only 1 location in history (not enough for velocity)
	history := []*domain.Location{
		{
			ID:            uuid.New(),
			ParticipantID: participantID,
			EntityID:      entityID,
			Latitude:      -23.500520,
			Longitude:     -46.633308,
			Timestamp:     now,
		},
	}

	mockRepo.On("GetLatestByParticipant", ctx, participantID, entityID).Return(latestLocation, nil)
	mockRepo.On("GetHistory", ctx, participantID, entityID, mock.Anything, mock.Anything).Return(history, nil)

	result, err := service.CalculateETA(ctx, participantID, entityID, targetLat, targetLng)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "simple", result.Method)
	mockRepo.AssertExpectations(t)
}

func TestETAService_CalculateETA_NoLatestLocation(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	cfg := &config.OSRMConfig{Enabled: false}
	service := eta.NewETAService(mockRepo, cfg)
	ctx := context.Background()

	participantID := uuid.New()
	entityID := uuid.New()

	mockRepo.On("GetLatestByParticipant", ctx, participantID, entityID).Return(nil, nil)

	result, err := service.CalculateETA(ctx, participantID, entityID, -23.550520, -46.633308)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no location data available")
	mockRepo.AssertExpectations(t)
}

func TestETAService_CalculateETA_RepositoryError(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	cfg := &config.OSRMConfig{Enabled: false}
	service := eta.NewETAService(mockRepo, cfg)
	ctx := context.Background()

	participantID := uuid.New()
	entityID := uuid.New()

	mockRepo.On("GetLatestByParticipant", ctx, participantID, entityID).Return(nil, errors.New("database error"))

	result, err := service.CalculateETA(ctx, participantID, entityID, -23.550520, -46.633308)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get latest location")
	mockRepo.AssertExpectations(t)
}

func TestETAService_CalculateMultipleETAs(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	cfg := &config.OSRMConfig{Enabled: false}
	service := eta.NewETAService(mockRepo, cfg)
	ctx := context.Background()

	participant1 := uuid.New()
	participant2 := uuid.New()
	participant3 := uuid.New()
	entityID := uuid.New()
	now := time.Now()

	// Participant 1 has location
	loc1 := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: participant1,
		EntityID:      entityID,
		Latitude:      -23.500520,
		Longitude:     -46.633308,
		Timestamp:     now,
	}

	// Participant 2 has location
	loc2 := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: participant2,
		EntityID:      entityID,
		Latitude:      -23.510520,
		Longitude:     -46.633308,
		Timestamp:     now,
	}

	mockRepo.On("GetLatestByParticipant", ctx, participant1, entityID).Return(loc1, nil)
	mockRepo.On("GetLatestByParticipant", ctx, participant2, entityID).Return(loc2, nil)
	mockRepo.On("GetLatestByParticipant", ctx, participant3, entityID).Return(nil, nil) // No location

	mockRepo.On("GetHistory", ctx, participant1, entityID, mock.Anything, mock.Anything).Return([]*domain.Location{}, nil)
	mockRepo.On("GetHistory", ctx, participant2, entityID, mock.Anything, mock.Anything).Return([]*domain.Location{}, nil)

	participantIDs := []uuid.UUID{participant1, participant2, participant3}
	results, err := service.CalculateMultipleETAs(ctx, participantIDs, entityID, -23.550520, -46.633308)

	assert.NoError(t, err)
	assert.Len(t, results, 2) // Only 2 participants have locations
	mockRepo.AssertExpectations(t)
}

func TestETAResult(t *testing.T) {
	participantID := uuid.New()
	now := time.Now()

	result := eta.ETAResult{
		ParticipantID:  participantID,
		DistanceMeters: 5000.0,
		ETAMinutes:     15,
		Method:         "velocity",
		LastUpdate:     now,
	}

	assert.Equal(t, participantID, result.ParticipantID)
	assert.Equal(t, 5000.0, result.DistanceMeters)
	assert.Equal(t, 15, result.ETAMinutes)
	assert.Equal(t, "velocity", result.Method)
	assert.Equal(t, now, result.LastUpdate)
}
