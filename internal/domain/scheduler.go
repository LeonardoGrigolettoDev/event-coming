package domain

import (
	"time"

	"github.com/google/uuid"
)

// SchedulerAction represents the type of scheduled action
type SchedulerAction string

const (
	SchedulerActionConfirmation SchedulerAction = "confirmation"
	SchedulerActionReminder     SchedulerAction = "reminder"
	SchedulerActionClosure      SchedulerAction = "closure"
	SchedulerActionLocation     SchedulerAction = "location"
)

// SchedulerStatus represents the status of a scheduler
type SchedulerStatus string

const (
	SchedulerStatusPending   SchedulerStatus = "pending"
	SchedulerStatusProcessed SchedulerStatus = "processed"
	SchedulerStatusFailed    SchedulerStatus = "failed"
	SchedulerStatusSkipped   SchedulerStatus = "skipped"
)

// Scheduler represents a scheduled task/action
type Scheduler struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	OrganizationID uuid.UUID       `json:"organization_id" db:"organization_id"`
	EventID        uuid.UUID       `json:"event_id" db:"event_id"`
	InstanceID     *uuid.UUID      `json:"instance_id,omitempty" db:"instance_id"`
	Action         SchedulerAction `json:"action" db:"action"`
	Status         SchedulerStatus `json:"status" db:"status"`
	ScheduledAt    time.Time       `json:"scheduled_at" db:"scheduled_at"`
	ProcessedAt    *time.Time      `json:"processed_at,omitempty" db:"processed_at"`
	Retries        int             `json:"retries" db:"retries"`
	MaxRetries     int             `json:"max_retries" db:"max_retries"`
	ErrorMessage   *string         `json:"error_message,omitempty" db:"error_message"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateSchedulerInput holds data for creating a scheduler
type CreateSchedulerInput struct {
	EventID     uuid.UUID              `json:"event_id" validate:"required"`
	InstanceID  *uuid.UUID             `json:"instance_id,omitempty"`
	Action      SchedulerAction        `json:"action" validate:"required,oneof=confirmation reminder closure location"`
	ScheduledAt time.Time              `json:"scheduled_at" validate:"required"`
	MaxRetries  int                    `json:"max_retries" validate:"min=0,max=10"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
