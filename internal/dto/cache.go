package dto

import (
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
)

// ==================== CACHE DATA ====================

// ParticipantLocationData representa a localização de um participante no cache
type ParticipantLocationData struct {
	ParticipantID   uuid.UUID  `json:"participant_id"`
	ParticipantName string     `json:"participant_name"`
	Latitude        float64    `json:"latitude"`
	Longitude       float64    `json:"longitude"`
	Accuracy        *float64   `json:"accuracy,omitempty"`
	Speed           *float64   `json:"speed,omitempty"`
	Heading         *float64   `json:"heading,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ETA             *time.Time `json:"eta,omitempty"`
	ETAMinutes      *int       `json:"eta_minutes,omitempty"`
}

// ParticipantConfirmationData representa a confirmação de um participante
type ParticipantConfirmationData struct {
	ParticipantID   uuid.UUID                `json:"participant_id"`
	ParticipantName string                   `json:"participant_name"`
	PhoneNumber     string                   `json:"phone_number"`
	Status          domain.ParticipantStatus `json:"status"`
	ConfirmedAt     *time.Time               `json:"confirmed_at,omitempty"`
	CheckedInAt     *time.Time               `json:"checked_in_at,omitempty"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

// ==================== RESPONSE ====================

// EventCacheResponse representa todos os dados em cache de um evento
type EventCacheResponse struct {
	OrganizationID uuid.UUID                     `json:"organization_id"`
	EventID        uuid.UUID                     `json:"event_id"`
	Locations      []ParticipantLocationData     `json:"locations"`
	Confirmations  []ParticipantConfirmationData `json:"confirmations"`
	TotalLocations int                           `json:"total_locations"`
	TotalConfirmed int                           `json:"total_confirmed"`
	TotalPending   int                           `json:"total_pending"`
	TotalDenied    int                           `json:"total_denied"`
	FetchedAt      time.Time                     `json:"fetched_at"`
}
