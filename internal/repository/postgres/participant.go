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

type participantRepository struct {
	db *gorm.DB
}

// NewParticipantRepository creates a new participant repository
func NewParticipantRepository(db *gorm.DB) repository.ParticipantRepository {
	return &participantRepository{db: db}
}

func (r *participantRepository) Create(ctx context.Context, participant *domain.Participant) error {
	if participant.ID == uuid.Nil {
		participant.ID = uuid.New()
	}

	result := r.db.WithContext(ctx).Create(participant)
	return result.Error
}

func (r *participantRepository) GetByID(ctx context.Context, id uuid.UUID, entityID uuid.UUID) (*domain.Participant, error) {
	var participant domain.Participant

	result := r.db.WithContext(ctx).
		Where("id = ? AND entity_id = ?", id, entityID).
		First(&participant)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return &participant, nil
}

func (r *participantRepository) Update(ctx context.Context, id uuid.UUID, entityID uuid.UUID, input *domain.UpdateParticipantInput) error {
	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.PhoneNumber != nil {
		updates["phone_number"] = *input.PhoneNumber
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.Metadata != nil {
		updates["metadata"] = input.Metadata
	}

	if len(updates) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Model(&domain.Participant{}).
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

func (r *participantRepository) Delete(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND entity_id = ?", id, entityID).
		Delete(&domain.Participant{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *participantRepository) ListByEvent(ctx context.Context, eventID uuid.UUID, entityID uuid.UUID, page, perPage int) ([]*domain.Participant, int64, error) {
	var participants []*domain.Participant
	var total int64

	offset := (page - 1) * perPage

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&domain.Participant{}).
		Where("event_id = ? AND entity_id = ?", eventID, entityID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Where("event_id = ? AND entity_id = ?", eventID, entityID).
		Order("name ASC").
		Offset(offset).
		Limit(perPage).
		Find(&participants).Error; err != nil {
		return nil, 0, err
	}

	return participants, total, nil
}

func (r *participantRepository) ListByEventInstance(ctx context.Context, instanceID uuid.UUID, entityID uuid.UUID, page, perPage int) ([]*domain.Participant, int64, error) {
	var participants []*domain.Participant
	var total int64

	offset := (page - 1) * perPage

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&domain.Participant{}).
		Where("instance_id = ? AND entity_id = ?", instanceID, entityID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Where("instance_id = ? AND entity_id = ?", instanceID, entityID).
		Order("name ASC").
		Offset(offset).
		Limit(perPage).
		Find(&participants).Error; err != nil {
		return nil, 0, err
	}

	return participants, total, nil
}

func (r *participantRepository) UpdateStatus(ctx context.Context, id uuid.UUID, entityID uuid.UUID, status domain.ParticipantStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// Set confirmed_at or checked_in_at based on status
	now := time.Now()
	switch status {
	case domain.ParticipantStatusConfirmed:
		updates["confirmed_at"] = now
	case domain.ParticipantStatusCheckedIn:
		updates["checked_in_at"] = now
	}

	result := r.db.WithContext(ctx).
		Model(&domain.Participant{}).
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

func (r *participantRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string, eventID uuid.UUID, entityID uuid.UUID) (*domain.Participant, error) {
	var participant domain.Participant

	result := r.db.WithContext(ctx).
		Where("phone_number = ? AND event_id = ? AND entity_id = ?", phoneNumber, eventID, entityID).
		First(&participant)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return &participant, nil
}

// GetActiveByPhoneNumber finds a participant by phone number in active events
// Returns the most recent participant with an active event
func (r *participantRepository) GetActiveByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Participant, error) {
	var participant domain.Participant

	// Join with events to find participants in active events
	result := r.db.WithContext(ctx).
		Joins("JOIN events ON events.id = participants.event_id").
		Where("participants.phone_number = ?", phoneNumber).
		Where("events.status = ?", domain.EventStatusActive).
		Where("events.start_time <= ? AND events.end_time >= ?", time.Now().Add(24*time.Hour), time.Now()).
		Order("events.start_time DESC").
		First(&participant)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return &participant, nil
}
