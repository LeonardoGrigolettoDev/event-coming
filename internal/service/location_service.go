package service

import (
	"context"
	"time"

	"event-coming/internal/cache"
	"event-coming/internal/domain"
	"event-coming/internal/dto"
	"event-coming/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LocationService handles location business logic
type LocationService struct {
	locationRepo    repository.LocationRepository
	participantRepo repository.ParticipantRepository
	eventRepo       repository.EventRepository
	locationBuffer  *cache.LocationBuffer
	logger          *zap.Logger
}

// NewLocationService creates a new location service
func NewLocationService(
	locationRepo repository.LocationRepository,
	participantRepo repository.ParticipantRepository,
	eventRepo repository.EventRepository,
	locationBuffer *cache.LocationBuffer,
	logger *zap.Logger,
) *LocationService {
	return &LocationService{
		locationRepo:    locationRepo,
		participantRepo: participantRepo,
		eventRepo:       eventRepo,
		locationBuffer:  locationBuffer,
		logger:          logger,
	}
}

// CreateLocation saves a new location for a participant
func (s *LocationService) CreateLocation(
	ctx context.Context,
	participantID uuid.UUID,
	entityID uuid.UUID,
	req *dto.CreateLocationRequest,
) (*dto.LocationResponse, error) {
	// Get participant to validate and get event info
	participant, err := s.participantRepo.GetByID(ctx, participantID, entityID)
	if err != nil {
		return nil, err
	}
	if participant == nil {
		return nil, domain.ErrNotFound
	}

	timestamp := time.Now()
	if req.Timestamp != nil {
		timestamp = *req.Timestamp
	}

	location := &domain.Location{
		ID:            uuid.New(),
		ParticipantID: participantID,
		EventID:       participant.EventID,
		EntityID:      entityID,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		Accuracy:      req.Accuracy,
		Altitude:      req.Altitude,
		Speed:         req.Speed,
		Heading:       req.Heading,
		Timestamp:     timestamp,
	}

	// Save to buffer (Redis) for real-time + batch processing
	if s.locationBuffer != nil {
		if err := s.locationBuffer.Push(ctx, location); err != nil {
			s.logger.Warn("Failed to push location to buffer", zap.Error(err))
			// Don't fail - still save to DB
		}
	}

	// Save to database
	if err := s.locationRepo.Create(ctx, location); err != nil {
		return nil, err
	}

	return dto.ToLocationResponse(location), nil
}

// GetLatestLocation gets the latest location for a participant
func (s *LocationService) GetLatestLocation(
	ctx context.Context,
	participantID uuid.UUID,
	entityID uuid.UUID,
) (*dto.LocationResponse, error) {
	location, err := s.locationRepo.GetLatestByParticipant(ctx, participantID, entityID)
	if err != nil {
		return nil, err
	}
	return dto.ToLocationResponse(location), nil
}

// GetLocationHistory gets location history for a participant
func (s *LocationService) GetLocationHistory(
	ctx context.Context,
	participantID uuid.UUID,
	entityID uuid.UUID,
	from, to time.Time,
) ([]*dto.LocationResponse, error) {
	locations, err := s.locationRepo.GetHistory(ctx, participantID, entityID, from, to)
	if err != nil {
		return nil, err
	}
	return dto.ToLocationResponseList(locations), nil
}

// GetEventLocations gets latest locations for all participants in an event
func (s *LocationService) GetEventLocations(
	ctx context.Context,
	eventID uuid.UUID,
	entityID uuid.UUID,
) ([]*dto.LocationResponse, error) {
	locations, err := s.locationRepo.GetLatestByEvent(ctx, eventID, entityID)
	if err != nil {
		return nil, err
	}
	return dto.ToLocationResponseList(locations), nil
}
