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
	SubscriptionPlanFree         SubscriptionPlan = "free"
	SubscriptionPlanBasic        SubscriptionPlan = "basic"
	SubscriptionPlanProfessional SubscriptionPlan = "professional"
	SubscriptionPlanEnterprise   SubscriptionPlan = "enterprise"
)

// Organization represents an organization/tenant
type Organization struct {
	ID               uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name             string           `json:"name" gorm:"type:varchar(100);not null"`
	Type             OrganizationType `json:"type" gorm:"type:varchar(50);not null"`
	SubscriptionPlan SubscriptionPlan `json:"subscription_plan" gorm:"column:subscription_plan;type:varchar(50);not null"`
	MaxEvents        int              `json:"max_events" gorm:"column:max_events;default:10"`
	MaxParticipants  int              `json:"max_participants" gorm:"column:max_participants;default:100"`
	Active           bool             `json:"active" gorm:"default:true"`
	CreatedAt        time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Organization
func (Organization) TableName() string {
	return "organizations"
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
