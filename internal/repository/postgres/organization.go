package postgres

import (
	"context"
	"errors"

	"event-coming/internal/domain"
	"event-coming/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(db *gorm.DB) repository.OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	if org.ID == uuid.Nil {
		org.ID = uuid.New()
	}

	result := r.db.WithContext(ctx).Create(org)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	var org domain.Organization

	result := r.db.WithContext(ctx).First(&org, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return &org, nil
}

func (r *organizationRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateOrganizationInput) error {
	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.SubscriptionPlan != nil {
		updates["subscription_plan"] = *input.SubscriptionPlan
	}
	if input.MaxEvents != nil {
		updates["max_events"] = *input.MaxEvents
	}
	if input.MaxParticipants != nil {
		updates["max_participants"] = *input.MaxParticipants
	}
	if input.Active != nil {
		updates["active"] = *input.Active
	}

	if len(updates) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Model(&domain.Organization{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *organizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Organization{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *organizationRepository) List(ctx context.Context, page, perPage int) ([]*domain.Organization, int64, error) {
	var orgs []*domain.Organization
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.Organization{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&orgs)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return orgs, total, nil
}
