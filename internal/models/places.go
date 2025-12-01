package models

import (
	"time"

	"github.com/google/uuid"
)

type Place struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	LocationDescription string    `json:"location_description" db:"location_description"`
	Longitude           float32   `json:"longitude" db:"longitude"`
	Latitude            float32   `json:"latitude" db:"latitude"`
	Altitude            float32   `json:"altitude" db:"altitude"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}
