package eta

import (
	"context"
	"fmt"
	"time"

	"event-coming/internal/config"
	"event-coming/internal/domain"
	"event-coming/internal/repository"

	"github.com/google/uuid"
)

// ETAService provides ETA calculation with multiple strategies
type ETAService struct {
	locationRepo  repository.LocationRepository
	velocityCalc  *VelocityCalculator
	osrmEnabled   bool
}

// NewETAService creates a new ETA service
func NewETAService(
	locationRepo repository.LocationRepository,
	cfg *config.OSRMConfig,
) *ETAService {
	return &ETAService{
		locationRepo:  locationRepo,
		velocityCalc:  NewVelocityCalculator(),
		osrmEnabled:   cfg.Enabled,
	}
}

// CalculateETA calculates ETA for a participant to reach an event location
func (s *ETAService) CalculateETA(
	ctx context.Context,
	participantID, orgID uuid.UUID,
	targetLat, targetLng float64,
) (*ETAResult, error) {
	// Get latest location
	latestLoc, err := s.locationRepo.GetLatestByParticipant(ctx, participantID, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest location: %w", err)
	}

	if latestLoc == nil {
		return nil, fmt.Errorf("no location data available")
	}

	// Calculate straight-line distance
	distance := CalculateHaversineDistance(
		latestLoc.Latitude, latestLoc.Longitude,
		targetLat, targetLng,
	)

	// Try OSRM first if enabled (placeholder for future implementation)
	if s.osrmEnabled {
		// TODO: Implement OSRM route calculation
	}

	// Fallback to velocity-based calculation
	history, err := s.locationRepo.GetHistory(
		ctx,
		participantID,
		orgID,
		time.Now().Add(-15*time.Minute),
		time.Now(),
	)

	var etaMinutes int
	var method string

	if err == nil && len(history) >= 2 {
		velocity := s.velocityCalc.CalculateVelocity(ctx, history)
		etaMinutes = s.velocityCalc.CalculateETA(distance, velocity)
		method = "velocity"
	} else {
		// Simple estimation: assume average speed of 30 km/h
		avgSpeedMPS := 30000.0 / 3600.0 // 30 km/h in m/s
		etaMinutes = s.velocityCalc.CalculateETA(distance, avgSpeedMPS)
		method = "simple"
	}

	return &ETAResult{
		ParticipantID:  participantID,
		DistanceMeters: distance,
		ETAMinutes:     etaMinutes,
		Method:         method,
		LastUpdate:     latestLoc.Timestamp,
	}, nil
}

// CalculateMultipleETAs calculates ETAs for multiple participants
func (s *ETAService) CalculateMultipleETAs(
	ctx context.Context,
	participantIDs []uuid.UUID,
	orgID uuid.UUID,
	targetLat, targetLng float64,
) ([]*ETAResult, error) {
	results := make([]*ETAResult, 0, len(participantIDs))

	for _, pid := range participantIDs {
		result, err := s.CalculateETA(ctx, pid, orgID, targetLat, targetLng)
		if err != nil {
			// Log error but continue
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

// ETAResult represents the result of an ETA calculation
type ETAResult struct {
	ParticipantID  uuid.UUID `json:"participant_id"`
	DistanceMeters float64   `json:"distance_meters"`
	ETAMinutes     int       `json:"eta_minutes"`
	Method         string    `json:"method"`
	LastUpdate     time.Time `json:"last_update"`
}
