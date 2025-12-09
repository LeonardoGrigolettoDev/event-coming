package domain

import (
	"time"

	"github.com/google/uuid"
)

// StatusResourceType represents the type of resource whose status changed
type StatusResourceType string

const (
	StatusResourceEvent       StatusResourceType = "event"
	StatusResourceParticipant StatusResourceType = "participant"
	StatusResourceScheduler   StatusResourceType = "scheduler"
)

// StatusHistory tracks status changes for events, participants, etc.
type StatusHistory struct {
	ID           uuid.UUID          `json:"id" db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ResourceType StatusResourceType `json:"resource_type" db:"resource_type" gorm:"size:50;not null;index"`
	ResourceID   uuid.UUID          `json:"resource_id" db:"resource_id" gorm:"type:uuid;not null;index"`
	EntityID     uuid.UUID          `json:"entity_id" db:"entity_id" gorm:"type:uuid;not null;index"`
	OldStatus    string             `json:"old_status" db:"old_status" gorm:"size:50"`
	NewStatus    string             `json:"new_status" db:"new_status" gorm:"size:50;not null"`
	ChangedBy    *uuid.UUID         `json:"changed_by,omitempty" db:"changed_by" gorm:"type:uuid"` // User or system (nil for auto)
	Reason       *string            `json:"reason,omitempty" db:"reason" gorm:"size:500"`
	Metadata     map[string]any     `json:"metadata,omitempty" db:"metadata" gorm:"type:jsonb"`
	CreatedAt    time.Time          `json:"created_at" db:"created_at" gorm:"autoCreateTime;index"`
}

func (StatusHistory) TableName() string {
	return "status_histories"
}

// StatusHistoryQuery holds query parameters for status history
type StatusHistoryQuery struct {
	ResourceType StatusResourceType
	ResourceID   uuid.UUID
	EntityID     uuid.UUID
	Page         int
	PerPage      int
}
