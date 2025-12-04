package domain

import (
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of event
type EventType string

const (
	EventTypeDemand     EventType = "demand"     // Non-recurring events
	EventTypePeriodic   EventType = "periodic"   // Recurring events
)

// EventStatus represents the status of an event
type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusScheduled EventStatus = "scheduled"
	EventStatusActive    EventStatus = "active"
	EventStatusCompleted EventStatus = "completed"
	EventStatusCancelled EventStatus = "cancelled"
)

// Event represents an event
type Event struct {
	ID                uuid.UUID   `json:"id" db:"id"`
	OrganizationID    uuid.UUID   `json:"organization_id" db:"organization_id"`
	Name              string      `json:"name" db:"name"`
	Description       *string     `json:"description,omitempty" db:"description"`
	Type              EventType   `json:"type" db:"type"`
	Status            EventStatus `json:"status" db:"status"`
	LocationLat       float64     `json:"location_lat" db:"location_lat"`
	LocationLng       float64     `json:"location_lng" db:"location_lng"`
	LocationAddress   *string     `json:"location_address,omitempty" db:"location_address"`
	StartTime         time.Time   `json:"start_time" db:"start_time"`
	EndTime           *time.Time  `json:"end_time,omitempty" db:"end_time"`
	RRuleString       *string     `json:"rrule_string,omitempty" db:"rrule_string"` // For periodic events
	ConfirmationDeadline *time.Time `json:"confirmation_deadline,omitempty" db:"confirmation_deadline"`
	CreatedBy         uuid.UUID   `json:"created_by" db:"created_by"`
	CreatedAt         time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at" db:"updated_at"`
}

// EventInstance represents a specific instance of a recurring event
type EventInstance struct {
	ID             uuid.UUID   `json:"id" db:"id"`
	EventID        uuid.UUID   `json:"event_id" db:"event_id"`
	OrganizationID uuid.UUID   `json:"organization_id" db:"organization_id"`
	InstanceDate   time.Time   `json:"instance_date" db:"instance_date"`
	Status         EventStatus `json:"status" db:"status"`
	StartTime      time.Time   `json:"start_time" db:"start_time"`
	EndTime        *time.Time  `json:"end_time,omitempty" db:"end_time"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" db:"updated_at"`
}

// CreateEventInput holds data for creating an event
type CreateEventInput struct {
	Name                 string     `json:"name" validate:"required,min=3,max=200"`
	Description          *string    `json:"description,omitempty" validate:"omitempty,max=1000"`
	Type                 EventType  `json:"type" validate:"required,oneof=demand periodic"`
	LocationLat          float64    `json:"location_lat" validate:"required,latitude"`
	LocationLng          float64    `json:"location_lng" validate:"required,longitude"`
	LocationAddress      *string    `json:"location_address,omitempty" validate:"omitempty,max=500"`
	StartTime            time.Time  `json:"start_time" validate:"required"`
	EndTime              *time.Time `json:"end_time,omitempty"`
	RRuleString          *string    `json:"rrule_string,omitempty" validate:"omitempty,max=500"`
	ConfirmationDeadline *time.Time `json:"confirmation_deadline,omitempty"`
}

// UpdateEventInput holds data for updating an event
type UpdateEventInput struct {
	Name                 *string     `json:"name,omitempty" validate:"omitempty,min=3,max=200"`
	Description          *string     `json:"description,omitempty" validate:"omitempty,max=1000"`
	Status               *EventStatus `json:"status,omitempty" validate:"omitempty,oneof=draft scheduled active completed cancelled"`
	LocationLat          *float64    `json:"location_lat,omitempty" validate:"omitempty,latitude"`
	LocationLng          *float64    `json:"location_lng,omitempty" validate:"omitempty,longitude"`
	LocationAddress      *string     `json:"location_address,omitempty" validate:"omitempty,max=500"`
	StartTime            *time.Time  `json:"start_time,omitempty"`
	EndTime              *time.Time  `json:"end_time,omitempty"`
	ConfirmationDeadline *time.Time  `json:"confirmation_deadline,omitempty"`
}
