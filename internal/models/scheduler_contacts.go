package models

import (
	"time"

	"github.com/google/uuid"
)

type SchedulerContact struct {
	ID          uuid.UUID `json:"id" db:"id"`
	SchedulerID uuid.UUID `json:"scheduler_id" db:"scheduler_id"`
	ContactID   uuid.UUID `json:"contact_id" db:"contact_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
