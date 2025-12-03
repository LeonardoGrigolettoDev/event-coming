package domain

import (
	"time"

	"github.com/google/uuid"
)

// Location represents a geolocation point
type Location struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ParticipantID  uuid.UUID `json:"participant_id" db:"participant_id"`
	EventID        uuid.UUID `json:"event_id" db:"event_id"`
	InstanceID     *uuid.UUID `json:"instance_id,omitempty" db:"instance_id"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	Latitude       float64   `json:"latitude" db:"latitude"`
	Longitude      float64   `json:"longitude" db:"longitude"`
	Accuracy       *float64  `json:"accuracy,omitempty" db:"accuracy"`
	Altitude       *float64  `json:"altitude,omitempty" db:"altitude"`
	Speed          *float64  `json:"speed,omitempty" db:"speed"`
	Heading        *float64  `json:"heading,omitempty" db:"heading"`
	Timestamp      time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// CreateLocationInput holds data for creating a location
type CreateLocationInput struct {
	ParticipantID uuid.UUID  `json:"participant_id" validate:"required"`
	EventID       uuid.UUID  `json:"event_id" validate:"required"`
	InstanceID    *uuid.UUID `json:"instance_id,omitempty"`
	Latitude      float64    `json:"latitude" validate:"required,latitude"`
	Longitude     float64    `json:"longitude" validate:"required,longitude"`
	Accuracy      *float64   `json:"accuracy,omitempty" validate:"omitempty,min=0"`
	Altitude      *float64   `json:"altitude,omitempty"`
	Speed         *float64   `json:"speed,omitempty" validate:"omitempty,min=0"`
	Heading       *float64   `json:"heading,omitempty" validate:"omitempty,min=0,max=360"`
	Timestamp     *time.Time `json:"timestamp,omitempty"`
}
