package domain

import (
	"time"

	"github.com/google/uuid"
)

// EntityType representa o tipo da entidade
type EntityType string
type EntityPermission string
type DocumentType string
type EntityRelationship string

const (
	EntityTypeNaturalPerson EntityType = "natural person" // Pessoa física
	EntityTypeLegalEntity   EntityType = "legal entity"   // Pessoa jurídica
)

const (
	EntityPermissionAdmin       EntityPermission = "admin"
	EntityPermissionStakeholder EntityPermission = "stakeholder"
	EntityPermissionParticipant EntityPermission = "participant"
)

const (
	DocumentTypeCPF  = "cpf"
	DocumentTypeCNPJ = "cnpj"
	DocumentTypeRG   = "rg"
)

const (
	RelationshipParent    EntityRelationship = "parent"
	RelationshipChild     EntityRelationship = "child"
	RelationshipSpouse    EntityRelationship = "spouse"
	RelationshipGuardian  EntityRelationship = "guardian"
	RelationshipDependent EntityRelationship = "dependent"
	RelationshipEmployee  EntityRelationship = "employee"
	RelationshipManager   EntityRelationship = "manager"
)

type Entity struct {
	ID               uuid.UUID              `json:"id" db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Relationship     EntityRelationship     `json:"relationship,omitempty" db:"relationship" gorm:"size:50"`
	ParentID         *uuid.UUID             `json:"parent_id,omitempty" db:"parent_id" gorm:"type:uuid;index"` // Entidade pai (hierarquia)
	Type             EntityType             `json:"type" db:"type" gorm:"size:50;not null;default:'natural person';index"`
	Name             string                 `json:"name" db:"name" gorm:"size:200"`
	Email            *string                `json:"email,omitempty" db:"email" gorm:"size:255;index"`
	PhoneNumber      *string                `json:"phone_number,omitempty" db:"phone_number" gorm:"size:20;index"`
	Document         *string                `json:"document,omitempty" db:"document" gorm:"size:50;index"` // CPF, CNPJ, etc.
	IsActive         bool                   `json:"is_active" db:"is_active" gorm:"default:true"`
	Metadata         map[string]interface{} `json:"metadata,omitempty" db:"metadata" gorm:"type:jsonb"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	EntityPermission EntityPermission       `json:"entity_permission" db:"entity_permission" gorm:"size:50;not null;default:'Participant'"`
	DocumentType     DocumentType           `json:"document_type" db:"document_type" gorm:"size:20"`
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
