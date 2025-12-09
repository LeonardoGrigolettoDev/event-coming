package postgres

import (
	"context"

	"event-coming/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type statusHistoryRepository struct {
	db *gorm.DB
}

// NewStatusHistoryRepository creates a new status history repository
func NewStatusHistoryRepository(db *gorm.DB) *statusHistoryRepository {
	return &statusHistoryRepository{db: db}
}

// Create saves a new status history entry
func (r *statusHistoryRepository) Create(ctx context.Context, history *domain.StatusHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// ListByResource returns status history for a specific resource
func (r *statusHistoryRepository) ListByResource(
	ctx context.Context,
	resourceType domain.StatusResourceType,
	resourceID uuid.UUID,
	page, perPage int,
) ([]*domain.StatusHistory, int64, error) {
	var histories []*domain.StatusHistory
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.StatusHistory{}).
		Where("resource_type = ? AND resource_id = ?", resourceType, resourceID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// ListByEntity returns status history for all resources in an entity
func (r *statusHistoryRepository) ListByEntity(
	ctx context.Context,
	entityID uuid.UUID,
	resourceType *domain.StatusResourceType,
	page, perPage int,
) ([]*domain.StatusHistory, int64, error) {
	var histories []*domain.StatusHistory
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.StatusHistory{}).
		Where("entity_id = ?", entityID)

	if resourceType != nil {
		query = query.Where("resource_type = ?", *resourceType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}
