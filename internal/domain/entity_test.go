package domain_test

import (
	"testing"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEntity_TableName(t *testing.T) {
	entity := domain.Entity{}
	assert.Equal(t, "entities", entity.TableName())
}

func TestEntity_CanCreateEvents(t *testing.T) {
	tests := []struct {
		name      string
		entity    domain.Entity
		canCreate bool
	}{
		{
			name: "Active admin can create events",
			entity: domain.Entity{
				Active:           true,
				EntityPermission: domain.EntityPermissionAdmin,
			},
			canCreate: true,
		},
		{
			name: "Active stakeholder can create events",
			entity: domain.Entity{
				Active:           true,
				EntityPermission: domain.EntityPermissionStakeholder,
			},
			canCreate: true,
		},
		{
			name: "Active participant cannot create events",
			entity: domain.Entity{
				Active:           true,
				EntityPermission: domain.EntityPermissionParticipant,
			},
			canCreate: false,
		},
		{
			name: "Inactive admin cannot create events",
			entity: domain.Entity{
				Active:           false,
				EntityPermission: domain.EntityPermissionAdmin,
			},
			canCreate: false,
		},
		{
			name: "Inactive stakeholder cannot create events",
			entity: domain.Entity{
				Active:           false,
				EntityPermission: domain.EntityPermissionStakeholder,
			},
			canCreate: false,
		},
		{
			name: "Inactive participant cannot create events",
			entity: domain.Entity{
				Active:           false,
				EntityPermission: domain.EntityPermissionParticipant,
			},
			canCreate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.entity.CanCreateEvents()
			assert.Equal(t, tt.canCreate, result)
		})
	}
}

func TestEntityType_Constants(t *testing.T) {
	assert.Equal(t, domain.EntityType("natural person"), domain.EntityTypeNaturalPerson)
	assert.Equal(t, domain.EntityType("legal entity"), domain.EntityTypeLegalEntity)
}

func TestEntityPermission_Constants(t *testing.T) {
	assert.Equal(t, domain.EntityPermission("admin"), domain.EntityPermissionAdmin)
	assert.Equal(t, domain.EntityPermission("stakeholder"), domain.EntityPermissionStakeholder)
	assert.Equal(t, domain.EntityPermission("participant"), domain.EntityPermissionParticipant)
}

func TestDocumentType_Constants(t *testing.T) {
	assert.Equal(t, "cpf", domain.DocumentTypeCPF)
	assert.Equal(t, "cnpj", domain.DocumentTypeCNPJ)
	assert.Equal(t, "rg", domain.DocumentTypeRG)
}

func TestEntityRelationship_Constants(t *testing.T) {
	assert.Equal(t, domain.EntityRelationship("parent"), domain.RelationshipParent)
	assert.Equal(t, domain.EntityRelationship("child"), domain.RelationshipChild)
	assert.Equal(t, domain.EntityRelationship("spouse"), domain.RelationshipSpouse)
	assert.Equal(t, domain.EntityRelationship("guardian"), domain.RelationshipGuardian)
	assert.Equal(t, domain.EntityRelationship("dependent"), domain.RelationshipDependent)
	assert.Equal(t, domain.EntityRelationship("employee"), domain.RelationshipEmployee)
	assert.Equal(t, domain.EntityRelationship("manager"), domain.RelationshipManager)
}

func TestEntity_Fields(t *testing.T) {
	id := uuid.New()
	parentID := uuid.New()
	now := time.Now()
	email := "test@example.com"
	phone := "+5511999999999"
	doc := "12345678901"

	entity := domain.Entity{
		ID:               id,
		ParentID:         &parentID,
		Relationship:     domain.RelationshipChild,
		Type:             domain.EntityTypeNaturalPerson,
		Name:             "Test Entity",
		Email:            &email,
		PhoneNumber:      &phone,
		Document:         &doc,
		Active:           true,
		Metadata:         map[string]interface{}{"key": "value"},
		CreatedAt:        now,
		UpdatedAt:        now,
		EntityPermission: domain.EntityPermissionAdmin,
		DocumentType:     domain.DocumentTypeCPF,
	}

	assert.Equal(t, id, entity.ID)
	assert.Equal(t, &parentID, entity.ParentID)
	assert.Equal(t, domain.RelationshipChild, entity.Relationship)
	assert.Equal(t, domain.EntityTypeNaturalPerson, entity.Type)
	assert.Equal(t, "Test Entity", entity.Name)
	assert.Equal(t, &email, entity.Email)
	assert.Equal(t, &phone, entity.PhoneNumber)
	assert.Equal(t, &doc, entity.Document)
	assert.True(t, entity.Active)
	assert.Equal(t, "value", entity.Metadata["key"])
	assert.Equal(t, now, entity.CreatedAt)
	assert.Equal(t, now, entity.UpdatedAt)
	assert.Equal(t, domain.EntityPermissionAdmin, entity.EntityPermission)
	assert.Equal(t, domain.DocumentType(domain.DocumentTypeCPF), entity.DocumentType)
}

func TestCreateEntityInput(t *testing.T) {
	parentID := uuid.New()
	email := "test@example.com"
	phone := "+5511999999999"
	doc := "12345678901"

	input := domain.CreateEntityInput{
		ParentID:    &parentID,
		Type:        domain.EntityTypeLegalEntity,
		Name:        "Test Company",
		Email:       &email,
		PhoneNumber: &phone,
		Document:    &doc,
		Metadata:    map[string]interface{}{"industry": "tech"},
	}

	assert.Equal(t, &parentID, input.ParentID)
	assert.Equal(t, domain.EntityTypeLegalEntity, input.Type)
	assert.Equal(t, "Test Company", input.Name)
	assert.Equal(t, &email, input.Email)
	assert.Equal(t, &phone, input.PhoneNumber)
	assert.Equal(t, &doc, input.Document)
	assert.Equal(t, "tech", input.Metadata["industry"])
}

func TestUpdateEntityInput(t *testing.T) {
	parentID := uuid.New()
	entityType := domain.EntityTypeLegalEntity
	name := "Updated Name"
	email := "updated@example.com"
	phone := "+5511888888888"
	doc := "98765432101"
	isActive := false

	input := domain.UpdateEntityInput{
		ParentID:    &parentID,
		Type:        &entityType,
		Name:        &name,
		Email:       &email,
		PhoneNumber: &phone,
		Document:    &doc,
		IsActive:    &isActive,
		Metadata:    map[string]interface{}{"updated": true},
	}

	assert.Equal(t, &parentID, input.ParentID)
	assert.Equal(t, &entityType, input.Type)
	assert.Equal(t, &name, input.Name)
	assert.Equal(t, &email, input.Email)
	assert.Equal(t, &phone, input.PhoneNumber)
	assert.Equal(t, &doc, input.Document)
	assert.Equal(t, &isActive, input.IsActive)
	assert.True(t, input.Metadata["updated"].(bool))
}
