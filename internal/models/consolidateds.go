package models

import (
	"time"

	"github.com/google/uuid"
)

type Consolidated struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	SchedulerID  uuid.UUID              `json:"scheduler_id" db:"scheduler_id"`
	ContactID    uuid.UUID              `json:"contact_id" db:"contact_id"`
	FinalStatus  string                 `json:"final_status" db:"final_status"`
	Timeline     map[string]interface{} `json:"timeline" db:"timeline"`
	PayloadFinal map[string]interface{} `json:"payload_final" db:"payload_final"`
	StartedAt    *time.Time             `json:"started_at" db:"started_at"`
	FinishedAt   *time.Time             `json:"finished_at" db:"finished_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
}
