package domain

import (
	"time"

	"github.com/google/uuid"
)

// EntityType representa o tipo da entidade
type EntityType string
type EntityPermission string

const (
	EntityTypeIndividual EntityType = "individual" // Pessoa física
	EntityTypeCompany    EntityType = "company"    // Pessoa jurídica
)

const (
	EntityPermissionAdmin       EntityPermission = "Admin"
	EntityPermissionStakeholder EntityPermission = "Stakeholder"
	EntityPermissionParticipant EntityPermission = "Participant"
)

// Entity representa uma entidade cadastrada no sistema
// Pode ser uma pessoa física, empresa ou organização
// Organizações podem criar eventos, pessoas podem ser participants
type Entity struct {
	ID               uuid.UUID              `json:"id" db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ParentID         *uuid.UUID             `json:"parent_id,omitempty" db:"parent_id" gorm:"type:uuid;index"` // Entidade pai (hierarquia)
	Type             EntityType             `json:"type" db:"type" gorm:"size:50;not null;default:'individual';index"`
	Name             string                 `json:"name" db:"name" gorm:"size:200;not null"`
	Email            *string                `json:"email,omitempty" db:"email" gorm:"size:255;index"`
	PhoneNumber      *string                `json:"phone_number,omitempty" db:"phone_number" gorm:"size:20;index"`
	Document         *string                `json:"document,omitempty" db:"document" gorm:"size:50;index"` // CPF, CNPJ, etc.
	IsActive         bool                   `json:"is_active" db:"is_active" gorm:"default:true"`
	Metadata         map[string]interface{} `json:"metadata,omitempty" db:"metadata" gorm:"type:jsonb"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	EntityPermission EntityPermission       `json:"entity_permission" db:"entity_permission" gorm:"size:50;not null;default:'Participant'"`
	// Relacionamentos
	Parent       *Entity       `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children     []Entity      `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Participants []Participant `json:"participants,omitempty" gorm:"foreignKey:EntityID"`
	Events       []Event       `json:"events,omitempty" gorm:"foreignKey:EntityID"` // Eventos criados por esta entidade
}

// TableName define o nome da tabela no banco
func (Entity) TableName() string {
	return "entities"
}

// CanCreateEvents retorna true se a entidade pode criar eventos
// Todas as entidades podem criar eventos
func (e *Entity) CanCreateEvents() bool {
	return e.IsActive && (e.EntityPermission == EntityPermissionAdmin || e.EntityPermission == EntityPermissionStakeholder)
}

// CreateEntityInput holds data for creating an entity
type CreateEntityInput struct {
	ParentID    *uuid.UUID
	Type        EntityType
	Name        string
	Email       *string
	PhoneNumber *string
	Document    *string
	Metadata    map[string]interface{}
}

// UpdateEntityInput holds data for updating an entity
type UpdateEntityInput struct {
	ParentID    *uuid.UUID
	Type        *EntityType
	Name        *string
	Email       *string
	PhoneNumber *string
	Document    *string
	IsActive    *bool
	Metadata    map[string]interface{}
}
