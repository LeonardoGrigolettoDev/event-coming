package dto

import (
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
)

// ==================== CREATE ====================

// ParticipantInput representa um participante a ser criado junto com o evento
type ParticipantInput struct {
	Name        string                 `json:"name" validate:"required,min=2,max=100"`
	PhoneNumber string                 `json:"phone_number" validate:"required"`
	Email       *string                `json:"email,omitempty" validate:"omitempty,email"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SchedulerConfig representa a configuração de schedulers a serem criados
type SchedulerConfig struct {
	SendConfirmation     bool       `json:"send_confirmation"`
	ConfirmationTime     *time.Time `json:"confirmation_time"`
	SendReminder         bool       `json:"send_reminder"`
	ReminderTime         *time.Time `json:"reminder_time"`
	ReminderBeforeHours  *int       `json:"reminder_before_hours"`
	TrackLocation        bool       `json:"track_location"`
	LocationTrackingTime *time.Time `json:"location_tracking_time"`
}

// CreateEventRequest representa o request de criação de evento
type CreateEventRequest struct {
	Name                 string             `json:"name" validate:"required,min=3,max=200"`
	Description          *string            `json:"description,omitempty" validate:"omitempty,max=1000"`
	Type                 domain.EventType   `json:"type" validate:"required,oneof=demand periodic"`
	LocationLat          float64            `json:"location_lat" validate:"required"`
	LocationLng          float64            `json:"location_lng" validate:"required"`
	LocationAddress      *string            `json:"location_address,omitempty" validate:"omitempty,max=500"`
	StartTime            time.Time          `json:"start_time" validate:"required"`
	EndTime              *time.Time         `json:"end_time,omitempty"`
	RRuleString          *string            `json:"rrule_string,omitempty" validate:"omitempty,max=500"`
	ConfirmationDeadline *time.Time         `json:"confirmation_deadline,omitempty"`
	Participants         []ParticipantInput `json:"participants,omitempty" validate:"omitempty,max=100,dive"`
	Scheduler            *SchedulerConfig   `json:"scheduler,omitempty"`
}

// ==================== UPDATE ====================

// UpdateEventRequest representa o request de atualização
type UpdateEventRequest struct {
	Name                 *string             `json:"name,omitempty" validate:"omitempty,min=3,max=200"`
	Description          *string             `json:"description,omitempty" validate:"omitempty,max=1000"`
	Status               *domain.EventStatus `json:"status,omitempty"`
	LocationLat          *float64            `json:"location_lat,omitempty"`
	LocationLng          *float64            `json:"location_lng,omitempty"`
	LocationAddress      *string             `json:"location_address,omitempty" validate:"omitempty,max=500"`
	StartTime            *time.Time          `json:"start_time,omitempty"`
	EndTime              *time.Time          `json:"end_time,omitempty"`
	ConfirmationDeadline *time.Time          `json:"confirmation_deadline,omitempty"`
}

// ==================== RESPONSE ====================

// EventResponse representa a resposta com dados do evento
type EventResponse struct {
	ID                   uuid.UUID              `json:"id"`
	EntityID             uuid.UUID              `json:"entity_id"`
	Name                 string                 `json:"name"`
	Description          *string                `json:"description,omitempty"`
	Type                 domain.EventType       `json:"type"`
	Status               domain.EventStatus     `json:"status"`
	LocationLat          float64                `json:"location_lat"`
	LocationLng          float64                `json:"location_lng"`
	LocationAddress      *string                `json:"location_address,omitempty"`
	StartTime            time.Time              `json:"start_time"`
	EndTime              *time.Time             `json:"end_time,omitempty"`
	RRuleString          *string                `json:"rrule_string,omitempty"`
	ConfirmationDeadline *time.Time             `json:"confirmation_deadline,omitempty"`
	CreatedBy            uuid.UUID              `json:"created_by"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	Participants         []*ParticipantResponse `json:"participants,omitempty"`
	SchedulersCreated    int                    `json:"schedulers_created,omitempty"`
}

// ToEventResponse converte domain.Event para EventResponse
func ToEventResponse(e *domain.Event) *EventResponse {
	return &EventResponse{
		ID:                   e.ID,
		EntityID:             e.EntityID,
		Name:                 e.Name,
		Description:          e.Description,
		Type:                 e.Type,
		Status:               e.Status,
		LocationLat:          e.LocationLat,
		LocationLng:          e.LocationLng,
		LocationAddress:      e.LocationAddress,
		StartTime:            e.StartTime,
		EndTime:              e.EndTime,
		RRuleString:          e.RRuleString,
		ConfirmationDeadline: e.ConfirmationDeadline,
		CreatedBy:            e.CreatedBy,
		CreatedAt:            e.CreatedAt,
		UpdatedAt:            e.UpdatedAt,
	}
}
