package postgres

import (
	"context"
	"errors"

	"event-coming/internal/domain"
	"event-coming/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *gorm.DB) repository.EventRepository {
	return &eventRepository{db: db}
}

// ==================== EVENT CRUD ====================

func (r *eventRepository) Create(ctx context.Context, event *domain.Event) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}

	result := r.db.WithContext(ctx).Create(event)
	return result.Error
}

func (r *eventRepository) GetByID(ctx context.Context, id uuid.UUID, entityID uuid.UUID) (*domain.Event, error) {
	var event domain.Event

	result := r.db.WithContext(ctx).
		Where("id = ? AND entity_id = ?", id, entityID).
		First(&event)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return &event, nil
}

func (r *eventRepository) Update(ctx context.Context, id uuid.UUID, entityID uuid.UUID, input *domain.UpdateEventInput) error {
	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.LocationLat != nil {
		updates["location_lat"] = *input.LocationLat
	}
	if input.LocationLng != nil {
		updates["location_lng"] = *input.LocationLng
	}
	if input.LocationAddress != nil {
		updates["location_address"] = *input.LocationAddress
	}
	if input.StartTime != nil {
		updates["start_time"] = *input.StartTime
	}
	if input.EndTime != nil {
		updates["end_time"] = *input.EndTime
	}
	if input.ConfirmationDeadline != nil {
		updates["confirmation_deadline"] = *input.ConfirmationDeadline
	}

	if len(updates) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Model(&domain.Event{}).
		Where("id = ? AND entity_id = ?", id, entityID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *eventRepository) Delete(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND entity_id = ?", id, entityID).
		Delete(&domain.Event{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *eventRepository) List(ctx context.Context, entityID uuid.UUID, page, perPage int) ([]*domain.Event, int64, error) {
	var events []*domain.Event
	var total int64

	offset := (page - 1) * perPage

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&domain.Event{}).
		Where("entity_id = ?", entityID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Where("entity_id = ?", entityID).
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *eventRepository) ListByStatus(ctx context.Context, entityID uuid.UUID, status domain.EventStatus, page, perPage int) ([]*domain.Event, int64, error) {
	var events []*domain.Event
	var total int64

	offset := (page - 1) * perPage

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&domain.Event{}).
		Where("entity_id = ? AND status = ?", entityID, status).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Where("entity_id = ? AND status = ?", entityID, status).
		Order("start_time ASC").
		Offset(offset).
		Limit(perPage).
		Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// ==================== EVENT INSTANCE ====================

func (r *eventRepository) CreateInstance(ctx context.Context, instance *domain.EventInstance) error {
	if instance.ID == uuid.Nil {
		instance.ID = uuid.New()
	}

	result := r.db.WithContext(ctx).Create(instance)
	return result.Error
}

func (r *eventRepository) GetInstanceByID(ctx context.Context, id uuid.UUID, entityID uuid.UUID) (*domain.EventInstance, error) {
	var instance domain.EventInstance

	result := r.db.WithContext(ctx).
		Where("id = ? AND entity_id = ?", id, entityID).
		First(&instance)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return &instance, nil
}

func (r *eventRepository) ListInstances(ctx context.Context, eventID uuid.UUID, entityID uuid.UUID) ([]*domain.EventInstance, error) {
	var instances []*domain.EventInstance

	result := r.db.WithContext(ctx).
		Where("event_id = ? AND entity_id = ?", eventID, entityID).
		Order("instance_date ASC").
		Find(&instances)

	if result.Error != nil {
		return nil, result.Error
	}

	return instances, nil
}
