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

	// Get event to use endTime for cache TTL
	event, err := s.eventRepo.GetByID(ctx, participant.EventID, entityID)
	if err != nil {
		s.logger.Warn("Failed to get event for cache TTL", zap.Error(err))
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

	// Save to Redis cache with TTL based on event end time
	if s.locationBuffer != nil {
		if event != nil && event.EndTime != nil {
			// Use event end time for TTL
			if err := s.locationBuffer.SetLatestLocation(ctx, location, *event.EndTime); err != nil {
				s.logger.Warn("Failed to set latest location in cache", zap.Error(err))
			}
		} else {
			// Fallback to default 24h TTL
			if err := s.locationBuffer.Push(ctx, location); err != nil {
				s.logger.Warn("Failed to push location to buffer", zap.Error(err))
			}
		}
	}

	// Save to database
	if err := s.locationRepo.Create(ctx, location); err != nil {
		return nil, err
	}

	return dto.ToLocationResponse(location), nil
}

// GetLatestLocation gets the latest location for a participant
// First tries Redis cache, then falls back to database
func (s *LocationService) GetLatestLocation(
	ctx context.Context,
	participantID uuid.UUID,
	entityID uuid.UUID,
) (*dto.LocationResponse, error) {
	// Get participant to get eventID for cache key
	participant, err := s.participantRepo.GetByID(ctx, participantID, entityID)
	if err != nil {
		return nil, err
	}

	// Try to get from Redis cache first
	if s.locationBuffer != nil {
		cachedLoc, err := s.locationBuffer.GetLatestLocation(ctx, participant.EventID, participantID)
		if err != nil {
			s.logger.Warn("Failed to get location from cache", zap.Error(err))
		} else if cachedLoc != nil {
			return dto.ToLocationResponse(cachedLoc), nil
		}
	}

	// Fallback to database
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
// First tries Redis cache, then falls back to database
func (s *LocationService) GetEventLocations(
	ctx context.Context,
	eventID uuid.UUID,
	entityID uuid.UUID,
) ([]*dto.LocationResponse, error) {
	// Try to get participant IDs for this event to check cache
	if s.locationBuffer != nil {
		participants, _, err := s.participantRepo.ListByEvent(ctx, eventID, entityID, 1, 1000)
		if err == nil && len(participants) > 0 {
			participantIDs := make([]uuid.UUID, len(participants))
			for i, p := range participants {
				participantIDs[i] = p.ID
			}

			cachedLocations, err := s.locationBuffer.GetLatestLocationsForEvent(ctx, eventID, participantIDs)
			if err != nil {
				s.logger.Warn("Failed to get locations from cache", zap.Error(err))
			} else if len(cachedLocations) > 0 {
				return dto.ToLocationResponseList(cachedLocations), nil
			}
		}
	}

	// Fallback to database
	locations, err := s.locationRepo.GetLatestByEvent(ctx, eventID, entityID)
	if err != nil {
		return nil, err
	}
	return dto.ToLocationResponseList(locations), nil
}
