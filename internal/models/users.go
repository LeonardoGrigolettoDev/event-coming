package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Email          uuid.UUID `json:"scheduler_id" db:"scheduler_id"`
	ContactID      uuid.UUID `json:"contact_id" db:"contact_id"`
	HashedPassword string    `json:"hashed_password" db:"hashed_password"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
