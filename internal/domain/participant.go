package domain

import (
	"time"

	"github.com/google/uuid"
)

// ParticipantStatus represents the status of a participant
type ParticipantStatus string

const (
	ParticipantStatusPending   ParticipantStatus = "pending"
	ParticipantStatusConfirmed ParticipantStatus = "confirmed"
	ParticipantStatusDenied    ParticipantStatus = "denied"
	ParticipantStatusCheckedIn ParticipantStatus = "checked_in"
	ParticipantStatusNoShow    ParticipantStatus = "no_show"
)

// Participant represents a participant in an event
type Participant struct {
	ID          uuid.UUID              `json:"id" db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EventID     uuid.UUID              `json:"event_id" db:"event_id" gorm:"type:uuid;not null;index"`
	InstanceID  *uuid.UUID             `json:"instance_id,omitempty" db:"instance_id" gorm:"type:uuid;index"`
	EntityID    uuid.UUID              `json:"entity_id" db:"entity_id" gorm:"type:uuid;not null;index"`          // Entidade dona do evento
	RefEntityID *uuid.UUID             `json:"ref_entity_id,omitempty" db:"ref_entity_id" gorm:"type:uuid;index"` // Referência opcional para entidade cadastrada do participante
	Status      ParticipantStatus      `json:"status" db:"status" gorm:"size:50;not null;default:'pending'"`
	ConfirmedAt *time.Time             `json:"confirmed_at,omitempty" db:"confirmed_at"`
	CheckedInAt *time.Time             `json:"checked_in_at,omitempty" db:"checked_in_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" db:"metadata" gorm:"type:jsonb"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`

	// Relacionamento
	Entity    *Entity `json:"entity,omitempty" gorm:"foreignKey:EntityID"`
	RefEntity *Entity `json:"ref_entity,omitempty" gorm:"foreignKey:RefEntityID"`
}

func (Participant) TableName() string {
	return "participants"
}

// CreateParticipantInput holds data for creating a participant
type CreateParticipantInput struct {
	EventID     uuid.UUID              `json:"event_id" validate:"required"`
	InstanceID  *uuid.UUID             `json:"instance_id,omitempty"`
	Name        string                 `json:"name" validate:"required,min=2,max=100"`
	PhoneNumber string                 `json:"phone_number" validate:"required,e164"`
	Email       *string                `json:"email,omitempty" validate:"omitempty,email"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateParticipantInput holds data for updating a participant
type UpdateParticipantInput struct {
	Name        *string                `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	PhoneNumber *string                `json:"phone_number,omitempty" validate:"omitempty,e164"`
	Email       *string                `json:"email,omitempty" validate:"omitempty,email"`
	Status      *ParticipantStatus     `json:"status,omitempty" validate:"omitempty,oneof=pending confirmed denied checked_in no_show"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ParticipantDistance holds participant distance information
type ParticipantDistance struct {
	ParticipantID uuid.UUID `json:"participant_id"`
	Name          string    `json:"name"`
	Distance      float64   `json:"distance_meters"`
	ETA           *int      `json:"eta_minutes,omitempty"`
	LastUpdate    time.Time `json:"last_update"`
}
