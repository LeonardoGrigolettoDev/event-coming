package dto_test

import (
	"testing"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateEntityRequest(t *testing.T) {
	parentID := uuid.New()
	email := "entity@example.com"
	phone := "+5511999999999"
	doc := "12345678901"

	req := dto.CreateEntityRequest{
		ParentID:    &parentID,
		Type:        domain.EntityTypeNaturalPerson,
		Name:        "Test Entity",
		Email:       &email,
		PhoneNumber: &phone,
		Document:    &doc,
		Metadata:    map[string]interface{}{"key": "value"},
	}

	assert.Equal(t, &parentID, req.ParentID)
	assert.Equal(t, domain.EntityTypeNaturalPerson, req.Type)
	assert.Equal(t, "Test Entity", req.Name)
	assert.Equal(t, &email, req.Email)
	assert.Equal(t, &phone, req.PhoneNumber)
	assert.Equal(t, &doc, req.Document)
	assert.Equal(t, "value", req.Metadata["key"])
}

func TestUpdateEntityRequest(t *testing.T) {
	parentID := uuid.New()
	entityType := domain.EntityTypeLegalEntity
	name := "Updated Name"
	email := "updated@example.com"
	phone := "+5511888888888"
	doc := "98765432101"
	isActive := false

	req := dto.UpdateEntityRequest{
		ParentID:    &parentID,
		Type:        &entityType,
		Name:        &name,
		Email:       &email,
		PhoneNumber: &phone,
		Document:    &doc,
		IsActive:    &isActive,
		Metadata:    map[string]interface{}{"updated": true},
	}

	assert.Equal(t, &parentID, req.ParentID)
	assert.Equal(t, &entityType, req.Type)
	assert.Equal(t, &name, req.Name)
	assert.Equal(t, &email, req.Email)
	assert.Equal(t, &phone, req.PhoneNumber)
	assert.Equal(t, &doc, req.Document)
	assert.Equal(t, &isActive, req.IsActive)
	assert.True(t, req.Metadata["updated"].(bool))
}

func TestEntityResponse(t *testing.T) {
	id := uuid.New()
	parentID := uuid.New()
	email := "entity@example.com"
	phone := "+5511999999999"
	doc := "12345678901"
	now := time.Now()

	resp := dto.EntityResponse{
		ID:               id,
		ParentID:         &parentID,
		Type:             domain.EntityTypeNaturalPerson,
		Name:             "Test Entity",
		Email:            &email,
		PhoneNumber:      &phone,
		Document:         &doc,
		IsActive:         true,
		EntityPermission: domain.EntityPermissionAdmin,
		Metadata:         map[string]interface{}{"key": "value"},
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, &parentID, resp.ParentID)
	assert.Equal(t, domain.EntityTypeNaturalPerson, resp.Type)
	assert.Equal(t, "Test Entity", resp.Name)
	assert.Equal(t, &email, resp.Email)
	assert.Equal(t, &phone, resp.PhoneNumber)
	assert.Equal(t, &doc, resp.Document)
	assert.True(t, resp.IsActive)
	assert.Equal(t, domain.EntityPermissionAdmin, resp.EntityPermission)
	assert.Equal(t, "value", resp.Metadata["key"])
}

func TestEntityResponse_WithChildren(t *testing.T) {
	childID := uuid.New()
	parentID := uuid.New()
	now := time.Now()

	child := dto.EntityResponse{
		ID:       childID,
		ParentID: &parentID,
		Type:     domain.EntityTypeNaturalPerson,
		Name:     "Child Entity",
		IsActive: true,
	}

	parent := dto.EntityResponse{
		ID:        parentID,
		Type:      domain.EntityTypeLegalEntity,
		Name:      "Parent Entity",
		IsActive:  true,
		Children:  []*dto.EntityResponse{&child},
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Len(t, parent.Children, 1)
	assert.Equal(t, childID, parent.Children[0].ID)
	assert.Equal(t, "Child Entity", parent.Children[0].Name)
}

func TestToEntityResponse(t *testing.T) {
	id := uuid.New()
	email := "entity@example.com"
	now := time.Now()

	entity := &domain.Entity{
		ID:               id,
		Type:             domain.EntityTypeNaturalPerson,
		Name:             "Test Entity",
		Email:            &email,
		Active:           true,
		EntityPermission: domain.EntityPermissionAdmin,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	resp := dto.ToEntityResponse(entity)

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, domain.EntityTypeNaturalPerson, resp.Type)
	assert.Equal(t, "Test Entity", resp.Name)
	assert.Equal(t, &email, resp.Email)
	assert.True(t, resp.IsActive)
	assert.Equal(t, domain.EntityPermissionAdmin, resp.EntityPermission)
}

func TestToEntityResponse_Nil(t *testing.T) {
	resp := dto.ToEntityResponse(nil)
	assert.Nil(t, resp)
}

func TestToEntityResponse_WithChildren(t *testing.T) {
	parentID := uuid.New()
	childID := uuid.New()
	now := time.Now()

	parent := &domain.Entity{
		ID:               parentID,
		Type:             domain.EntityTypeLegalEntity,
		Name:             "Parent Entity",
		Active:           true,
		EntityPermission: domain.EntityPermissionAdmin,
		CreatedAt:        now,
		UpdatedAt:        now,
		Children: []domain.Entity{
			{
				ID:               childID,
				ParentID:         &parentID,
				Type:             domain.EntityTypeNaturalPerson,
				Name:             "Child Entity",
				Active:           true,
				EntityPermission: domain.EntityPermissionParticipant,
			},
		},
	}

	resp := dto.ToEntityResponse(parent)

	assert.Equal(t, parentID, resp.ID)
	assert.Len(t, resp.Children, 1)
	assert.Equal(t, childID, resp.Children[0].ID)
	assert.Equal(t, "Child Entity", resp.Children[0].Name)
}

func TestToEntityResponseList(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	now := time.Now()

	entities := []*domain.Entity{
		{
			ID:               id1,
			Type:             domain.EntityTypeNaturalPerson,
			Name:             "Entity 1",
			Active:           true,
			EntityPermission: domain.EntityPermissionAdmin,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               id2,
			Type:             domain.EntityTypeLegalEntity,
			Name:             "Entity 2",
			Active:           true,
			EntityPermission: domain.EntityPermissionStakeholder,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
	}

	responses := dto.ToEntityResponseList(entities)

	assert.Len(t, responses, 2)
	assert.Equal(t, id1, responses[0].ID)
	assert.Equal(t, "Entity 1", responses[0].Name)
	assert.Equal(t, id2, responses[1].ID)
	assert.Equal(t, "Entity 2", responses[1].Name)
}

func TestToEntityResponseList_Empty(t *testing.T) {
	entities := []*domain.Entity{}
	responses := dto.ToEntityResponseList(entities)
	assert.Empty(t, responses)
}
