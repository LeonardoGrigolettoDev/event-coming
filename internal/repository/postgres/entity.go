package postgres

import (
	"context"
	"errors"

	"event-coming/internal/domain"
	"event-coming/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type entityRepository struct {
	db *gorm.DB
}

// NewEntityRepository creates a new Entity repository
func NewEntityRepository(db *gorm.DB) repository.EntityRepository {
	return &entityRepository{db: db}
}

// Create creates a new entity
func (r *entityRepository) Create(ctx context.Context, entity *domain.Entity) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// GetByID retrieves an entity by ID
func (r *entityRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Entity, error) {
	var entity domain.Entity
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// Update updates an entity
func (r *entityRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateEntityInput) error {
	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Type != nil {
		updates["type"] = *input.Type
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.PhoneNumber != nil {
		updates["phone_number"] = *input.PhoneNumber
	}
	if input.Document != nil {
		updates["document"] = *input.Document
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	if input.ParentID != nil {
		updates["parent_id"] = *input.ParentID
	}
	if input.Metadata != nil {
		updates["metadata"] = input.Metadata
	}

	if len(updates) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Model(&domain.Entity{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete deletes an entity
func (r *entityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.Entity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// List lists entities with pagination
func (r *entityRepository) List(ctx context.Context, page, perPage int) ([]*domain.Entity, int64, error) {
	var entities []*domain.Entity
	var total int64

	offset := (page - 1) * perPage

	if err := r.db.WithContext(ctx).
		Model(&domain.Entity{}).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// ListByParent lists entities by parent ID
func (r *entityRepository) ListByParent(ctx context.Context, parentID uuid.UUID, page, perPage int) ([]*domain.Entity, int64, error) {
	var entities []*domain.Entity
	var total int64

	offset := (page - 1) * perPage

	if err := r.db.WithContext(ctx).
		Model(&domain.Entity{}).
		Where("parent_id = ?", parentID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// GetByDocument retrieves an entity by document
func (r *entityRepository) GetByDocument(ctx context.Context, document string) (*domain.Entity, error) {
	var entity domain.Entity
	err := r.db.WithContext(ctx).
		Where("document = ?", document).
		First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}
