package dto

import (
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
)

// ==================== CREATE ====================

// CreateEntityRequest representa o request de criação de entidade
type CreateEntityRequest struct {
	ParentID    *uuid.UUID             `json:"parent_id,omitempty"`
	Type        domain.EntityType      `json:"type" validate:"required,oneof=individual company"`
	Name        string                 `json:"name" validate:"required,min=2,max=200"`
	Email       *string                `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber *string                `json:"phone_number,omitempty" validate:"omitempty,max=20"`
	Document    *string                `json:"document,omitempty" validate:"omitempty,max=50"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ==================== UPDATE ====================

// UpdateEntityRequest representa o request de atualização
type UpdateEntityRequest struct {
	ParentID    *uuid.UUID             `json:"parent_id,omitempty"`
	Type        *domain.EntityType     `json:"type,omitempty" validate:"omitempty,oneof=individual company"`
	Name        *string                `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Email       *string                `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber *string                `json:"phone_number,omitempty" validate:"omitempty,max=20"`
	Document    *string                `json:"document,omitempty" validate:"omitempty,max=50"`
	IsActive    *bool                  `json:"is_active,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ==================== RESPONSE ====================

// EntityResponse representa a resposta com dados da entidade
type EntityResponse struct {
	ID               uuid.UUID               `json:"id"`
	ParentID         *uuid.UUID              `json:"parent_id,omitempty"`
	Type             domain.EntityType       `json:"type"`
	Name             string                  `json:"name"`
	Email            *string                 `json:"email,omitempty"`
	PhoneNumber      *string                 `json:"phone_number,omitempty"`
	Document         *string                 `json:"document,omitempty"`
	IsActive         bool                    `json:"is_active"`
	EntityPermission domain.EntityPermission `json:"entity_permission"`
	Metadata         map[string]interface{}  `json:"metadata,omitempty"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	Children         []*EntityResponse       `json:"children,omitempty"`
}

// ToEntityResponse converte domain.Entity para EntityResponse
func ToEntityResponse(e *domain.Entity) *EntityResponse {
	if e == nil {
		return nil
	}

	resp := &EntityResponse{
		ID:               e.ID,
		ParentID:         e.ParentID,
		Type:             e.Type,
		Name:             e.Name,
		Email:            e.Email,
		PhoneNumber:      e.PhoneNumber,
		Document:         e.Document,
		IsActive:         e.Active,
		EntityPermission: e.EntityPermission,
		Metadata:         e.Metadata,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}

	// Converter children se existirem
	if len(e.Children) > 0 {
		resp.Children = make([]*EntityResponse, len(e.Children))
		for i, child := range e.Children {
			childCopy := child
			resp.Children[i] = ToEntityResponse(&childCopy)
		}
	}

	return resp
}

// ToEntityResponseList converte uma lista de domain.Entity para EntityResponse
func ToEntityResponseList(entities []*domain.Entity) []*EntityResponse {
	responses := make([]*EntityResponse, len(entities))
	for i, e := range entities {
		responses[i] = ToEntityResponse(e)
	}
	return responses
}
