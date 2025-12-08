package postgres

import (
	"context"
	"errors"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type schedulerRepository struct {
	db *gorm.DB
}

// NewSchedulerRepository creates a new scheduler repository
func NewSchedulerRepository(db *gorm.DB) repository.SchedulerRepository {
	return &schedulerRepository{db: db}
}

func (r *schedulerRepository) Create(ctx context.Context, scheduler *domain.Scheduler) error {
	if scheduler.ID == uuid.Nil {
		scheduler.ID = uuid.New()
	}

	result := r.db.WithContext(ctx).Create(scheduler)
	return result.Error
}

func (r *schedulerRepository) GetByID(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (*domain.Scheduler, error) {
	var scheduler domain.Scheduler

	result := r.db.WithContext(ctx).
		Where("id = ? AND organization_id = ?", id, orgID).
		First(&scheduler)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return &scheduler, nil
}

func (r *schedulerRepository) Update(ctx context.Context, scheduler *domain.Scheduler) error {
	result := r.db.WithContext(ctx).Save(scheduler)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *schedulerRepository) Delete(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND organization_id = ?", id, orgID).
		Delete(&domain.Scheduler{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *schedulerRepository) ListPending(ctx context.Context, before time.Time, limit int) ([]*domain.Scheduler, error) {
	var schedulers []*domain.Scheduler

	result := r.db.WithContext(ctx).
		Where("status = ? AND scheduled_at <= ? AND retries < max_retries", domain.SchedulerStatusPending, before).
		Order("scheduled_at ASC").
		Limit(limit).
		Find(&schedulers)

	if result.Error != nil {
		return nil, result.Error
	}

	return schedulers, nil
}

func (r *schedulerRepository) MarkAsProcessed(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error {
	now := time.Now()

	result := r.db.WithContext(ctx).
		Model(&domain.Scheduler{}).
		Where("id = ? AND organization_id = ?", id, orgID).
		Updates(map[string]interface{}{
			"status":       domain.SchedulerStatusProcessed,
			"processed_at": now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *schedulerRepository) MarkAsFailed(ctx context.Context, id uuid.UUID, orgID uuid.UUID, errorMsg string) error {
	result := r.db.WithContext(ctx).
		Model(&domain.Scheduler{}).
		Where("id = ? AND organization_id = ?", id, orgID).
		Updates(map[string]interface{}{
			"status":        domain.SchedulerStatusFailed,
			"error_message": errorMsg,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *schedulerRepository) IncrementRetries(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&domain.Scheduler{}).
		Where("id = ? AND organization_id = ?", id, orgID).
		UpdateColumn("retries", gorm.Expr("retries + 1"))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
