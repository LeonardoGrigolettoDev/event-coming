package dto

import (
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
)

// ==================== CREATE ====================

// CreateParticipantRequest representa o request de criação de participante
type CreateParticipantRequest struct {
	Name        string                 `json:"name" validate:"required,min=2,max=100"`
	PhoneNumber string                 `json:"phone_number" validate:"required"`
	Email       *string                `json:"email,omitempty" validate:"omitempty,email"`
	InstanceID  *uuid.UUID             `json:"instance_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// BatchCreateParticipantsRequest representa request de criação em lote
type BatchCreateParticipantsRequest struct {
	Participants []CreateParticipantRequest `json:"participants" validate:"required,min=1,max=100"`
}

// ==================== UPDATE ====================

// UpdateParticipantRequest representa o request de atualização
type UpdateParticipantRequest struct {
	Name        *string                   `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	PhoneNumber *string                   `json:"phone_number,omitempty"`
	Email       *string                   `json:"email,omitempty" validate:"omitempty,email"`
	Status      *domain.ParticipantStatus `json:"status,omitempty"`
	Metadata    map[string]interface{}    `json:"metadata,omitempty"`
}

// ==================== RESPONSE ====================

// ParticipantResponse representa a resposta com dados do participante
type ParticipantResponse struct {
	ID          uuid.UUID                `json:"id"`
	EventID     uuid.UUID                `json:"event_id"`
	InstanceID  *uuid.UUID               `json:"instance_id,omitempty"`
	EntityID    uuid.UUID                `json:"entity_id"`
	Name        string                   `json:"name"`
	PhoneNumber string                   `json:"phone_number"`
	Email       *string                  `json:"email,omitempty"`
	Status      domain.ParticipantStatus `json:"status"`
	ConfirmedAt *time.Time               `json:"confirmed_at,omitempty"`
	CheckedInAt *time.Time               `json:"checked_in_at,omitempty"`
	Metadata    map[string]interface{}   `json:"metadata,omitempty"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

// ToParticipantResponse converte domain.Participant para ParticipantResponse
func ToParticipantResponse(p *domain.Participant) *ParticipantResponse {
	return &ParticipantResponse{
		ID:          p.ID,
		EventID:     p.EventID,
		InstanceID:  p.InstanceID,
		EntityID:    p.EntityID,
		Name:        p.Name,
		PhoneNumber: p.PhoneNumber,
		Email:       p.Email,
		Status:      p.Status,
		ConfirmedAt: p.ConfirmedAt,
		CheckedInAt: p.CheckedInAt,
		Metadata:    p.Metadata,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
