package dto

import (
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/service/eta"

	"github.com/google/uuid"
)

// ==================== CREATE ====================

// CreateLocationRequest representa o request de envio de localização
type CreateLocationRequest struct {
	Latitude  float64    `json:"latitude" binding:"required"`
	Longitude float64    `json:"longitude" binding:"required"`
	Accuracy  *float64   `json:"accuracy,omitempty"`
	Altitude  *float64   `json:"altitude,omitempty"`
	Speed     *float64   `json:"speed,omitempty"`
	Heading   *float64   `json:"heading,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// ==================== RESPONSE ====================

// LocationResponse representa a resposta com dados de localização
type LocationResponse struct {
	ID            uuid.UUID `json:"id"`
	ParticipantID uuid.UUID `json:"participant_id"`
	EventID       uuid.UUID `json:"event_id"`
	EntityID      uuid.UUID `json:"entity_id"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Accuracy      *float64  `json:"accuracy,omitempty"`
	Altitude      *float64  `json:"altitude,omitempty"`
	Speed         *float64  `json:"speed,omitempty"`
	Heading       *float64  `json:"heading,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	CreatedAt     time.Time `json:"created_at"`
}

// ToLocationResponse converte domain.Location para LocationResponse
func ToLocationResponse(loc *domain.Location) *LocationResponse {
	if loc == nil {
		return nil
	}
	return &LocationResponse{
		ID:            loc.ID,
		ParticipantID: loc.ParticipantID,
		EventID:       loc.EventID,
		EntityID:      loc.EntityID,
		Latitude:      loc.Latitude,
		Longitude:     loc.Longitude,
		Accuracy:      loc.Accuracy,
		Altitude:      loc.Altitude,
		Speed:         loc.Speed,
		Heading:       loc.Heading,
		Timestamp:     loc.Timestamp,
		CreatedAt:     loc.CreatedAt,
	}
}

// ToLocationResponseList converte lista de locations
func ToLocationResponseList(locations []*domain.Location) []*LocationResponse {
	responses := make([]*LocationResponse, len(locations))
	for i, loc := range locations {
		responses[i] = ToLocationResponse(loc)
	}
	return responses
}

// ==================== ETA ====================

// ETAResponse representa a resposta de cálculo de ETA
type ETAResponse struct {
	ParticipantID  uuid.UUID `json:"participant_id"`
	DistanceMeters float64   `json:"distance_meters"`
	ETAMinutes     int       `json:"eta_minutes"`
	Method         string    `json:"method"`
	LastUpdate     time.Time `json:"last_update"`
}

// ToETAResponse converte eta.ETAResult para ETAResponse
func ToETAResponse(result *eta.ETAResult) *ETAResponse {
	if result == nil {
		return nil
	}
	return &ETAResponse{
		ParticipantID:  result.ParticipantID,
		DistanceMeters: result.DistanceMeters,
		ETAMinutes:     result.ETAMinutes,
		Method:         result.Method,
		LastUpdate:     result.LastUpdate,
	}
}

// ToETAResponseList converte lista de ETAResults
func ToETAResponseList(results []*eta.ETAResult) []*ETAResponse {
	responses := make([]*ETAResponse, len(results))
	for i, r := range results {
		responses[i] = ToETAResponse(r)
	}
	return responses
}

// EventETAResponse representa ETAs de todos participantes de um evento
type EventETAResponse struct {
	EventID      uuid.UUID      `json:"event_id"`
	EntityID     uuid.UUID      `json:"entity_id"`
	Participants []*ETAResponse `json:"participants"`
	FetchedAt    time.Time      `json:"fetched_at"`
}

// ToEventETAResponse cria EventETAResponse a partir de lista de ETAResults
func ToEventETAResponse(eventID, entityID uuid.UUID, results []*eta.ETAResult) *EventETAResponse {
	return &EventETAResponse{
		EventID:      eventID,
		EntityID:     entityID,
		Participants: ToETAResponseList(results),
		FetchedAt:    time.Now(),
	}
}
