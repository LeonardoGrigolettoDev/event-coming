package service

import (
	"context"

	"event-coming/internal/domain"
	"event-coming/internal/dto"
	"event-coming/internal/repository"

	"github.com/google/uuid"
)

// EntityService handles entity business logic
type EntityService struct {
	entityRepo repository.EntityRepository
}

// NewEntityService creates a new entity service
func NewEntityService(entityRepo repository.EntityRepository) *EntityService {
	return &EntityService{
		entityRepo: entityRepo,
	}
}

// Create creates a new entity
func (s *EntityService) Create(ctx context.Context, req *dto.CreateEntityRequest) (*dto.EntityResponse, error) {
	// Check if document already exists
	if req.Document != nil && *req.Document != "" {
		existing, err := s.entityRepo.GetByDocument(ctx, *req.Document)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, domain.ErrConflict
		}
	}

	// Validate parent if provided
	if req.ParentID != nil {
		parent, err := s.entityRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, domain.ErrNotFound
		}
	}

	entity := &domain.Entity{
		ID:          uuid.New(),
		ParentID:    req.ParentID,
		Type:        req.Type,
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Document:    req.Document,
		IsActive:    true,
		Metadata:    req.Metadata,
	}

	if err := s.entityRepo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return dto.ToEntityResponse(entity), nil
}

// GetByID retrieves an entity by ID
func (s *EntityService) GetByID(ctx context.Context, id uuid.UUID) (*dto.EntityResponse, error) {
	entity, err := s.entityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, domain.ErrNotFound
	}

	return dto.ToEntityResponse(entity), nil
}

// Update updates an entity
func (s *EntityService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateEntityRequest) (*dto.EntityResponse, error) {
	// Check if entity exists
	existing, err := s.entityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.ErrNotFound
	}

	// Check document uniqueness if updating
	if req.Document != nil && *req.Document != "" {
		docEntity, err := s.entityRepo.GetByDocument(ctx, *req.Document)
		if err != nil {
			return nil, err
		}
		if docEntity != nil && docEntity.ID != id {
			return nil, domain.ErrConflict
		}
	}

	// Validate parent if provided
	if req.ParentID != nil {
		if *req.ParentID == id {
			return nil, domain.ErrInvalidInput
		}
		parent, err := s.entityRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, domain.ErrNotFound
		}
	}

	input := &domain.UpdateEntityInput{
		ParentID:    req.ParentID,
		Type:        req.Type,
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Document:    req.Document,
		IsActive:    req.IsActive,
		Metadata:    req.Metadata,
	}

	if err := s.entityRepo.Update(ctx, id, input); err != nil {
		return nil, err
	}

	// Fetch updated entity
	updated, err := s.entityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToEntityResponse(updated), nil
}

// Delete deletes an entity
func (s *EntityService) Delete(ctx context.Context, id uuid.UUID) error {
	existing, err := s.entityRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrNotFound
	}

	return s.entityRepo.Delete(ctx, id)
}

// List lists entities with pagination
func (s *EntityService) List(ctx context.Context, page, perPage int) ([]*dto.EntityResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	entities, total, err := s.entityRepo.List(ctx, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	return dto.ToEntityResponseList(entities), total, nil
}

// ListByParent lists entities by parent ID
func (s *EntityService) ListByParent(ctx context.Context, parentID uuid.UUID, page, perPage int) ([]*dto.EntityResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	entities, total, err := s.entityRepo.ListByParent(ctx, parentID, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	return dto.ToEntityResponseList(entities), total, nil
}

// GetByDocument retrieves an entity by document
func (s *EntityService) GetByDocument(ctx context.Context, document string) (*dto.EntityResponse, error) {
	entity, err := s.entityRepo.GetByDocument(ctx, document)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, domain.ErrNotFound
	}

	return dto.ToEntityResponse(entity), nil
}
