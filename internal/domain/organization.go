package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrganizationType represents the type of organization
type OrganizationType string

const (
	OrganizationTypeSchool     OrganizationType = "school"
	OrganizationTypeEnterprise OrganizationType = "enterprise"
	OrganizationTypeEvent      OrganizationType = "event"
)

// SubscriptionPlan represents the subscription tier
type SubscriptionPlan string

const (
	SubscriptionPlanFree       SubscriptionPlan = "free"
	SubscriptionPlanBasic      SubscriptionPlan = "basic"
	SubscriptionPlanProfessional SubscriptionPlan = "professional"
	SubscriptionPlanEnterprise SubscriptionPlan = "enterprise"
)

// Organization represents an organization/tenant
type Organization struct {
	ID               uuid.UUID        `json:"id" db:"id"`
	Name             string           `json:"name" db:"name"`
	Type             OrganizationType `json:"type" db:"type"`
	SubscriptionPlan SubscriptionPlan `json:"subscription_plan" db:"subscription_plan"`
	MaxEvents        int              `json:"max_events" db:"max_events"`
	MaxParticipants  int              `json:"max_participants" db:"max_participants"`
	Active           bool             `json:"active" db:"active"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
}

// CreateOrganizationInput holds data for creating an organization
type CreateOrganizationInput struct {
	Name             string           `json:"name" validate:"required,min=3,max=100"`
	Type             OrganizationType `json:"type" validate:"required,oneof=school enterprise event"`
	SubscriptionPlan SubscriptionPlan `json:"subscription_plan" validate:"required,oneof=free basic professional enterprise"`
}

// UpdateOrganizationInput holds data for updating an organization
type UpdateOrganizationInput struct {
	Name             *string           `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	SubscriptionPlan *SubscriptionPlan `json:"subscription_plan,omitempty" validate:"omitempty,oneof=free basic professional enterprise"`
	MaxEvents        *int              `json:"max_events,omitempty" validate:"omitempty,min=0"`
	MaxParticipants  *int              `json:"max_participants,omitempty" validate:"omitempty,min=0"`
	Active           *bool             `json:"active,omitempty"`
}
