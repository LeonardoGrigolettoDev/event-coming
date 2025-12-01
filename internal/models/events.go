package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	SchedulerID uuid.UUID              `json:"scheduler_id" db:"scheduler_id"`
	ContactID   uuid.UUID              `json:"contact_id" db:"contact_id"`
	EventType   string                 `json:"event_type" db:"event_type"`
	Status      string                 `json:"status" db:"status"`
	Payload     map[string]interface{} `json:"payload" db:"payload"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}
