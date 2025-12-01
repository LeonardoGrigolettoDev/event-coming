package models

import (
	"time"

	"github.com/google/uuid"
)

type Scheduler struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	Description    string     `json:"description" db:"description"`
	ScheduleType   string     `json:"schedule_type" db:"schedule_type"`
	CronExpression *string    `json:"cron_expression" db:"cron_expression"`
	StartDate      *time.Time `json:"start_date" db:"start_date"`
	EndDate        *time.Time `json:"end_date" db:"end_date"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}
