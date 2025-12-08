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

type locationRepository struct {
	db *gorm.DB
}

// NewLocationRepository creates a new location repository
func NewLocationRepository(db *gorm.DB) repository.LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Create(ctx context.Context, location *domain.Location) error {
	if location.ID == uuid.Nil {
		location.ID = uuid.New()
	}

	result := r.db.WithContext(ctx).Create(location)
	return result.Error
}

func (r *locationRepository) BatchCreate(ctx context.Context, locations []*domain.Location) error {
	if len(locations) == 0 {
		return nil
	}

	// Assign UUIDs to locations without IDs
	for _, loc := range locations {
		if loc.ID == uuid.Nil {
			loc.ID = uuid.New()
		}
	}

	result := r.db.WithContext(ctx).CreateInBatches(locations, 100)
	return result.Error
}

func (r *locationRepository) GetLatestByParticipant(ctx context.Context, participantID uuid.UUID, entityID uuid.UUID) (*domain.Location, error) {
	var location domain.Location

	result := r.db.WithContext(ctx).
		Where("participant_id = ? AND entity_id = ?", participantID, entityID).
		Order("timestamp DESC").
		First(&location)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return &location, nil
}

func (r *locationRepository) GetLatestByEvent(ctx context.Context, eventID uuid.UUID, entityID uuid.UUID) ([]*domain.Location, error) {
	var locations []*domain.Location

	// Subquery to get latest location per participant
	subQuery := r.db.WithContext(ctx).
		Model(&domain.Location{}).
		Select("participant_id, MAX(timestamp) as max_timestamp").
		Where("event_id = ? AND entity_id = ?", eventID, entityID).
		Group("participant_id")

	result := r.db.WithContext(ctx).
		Where("event_id = ? AND entity_id = ?", eventID, entityID).
		Where("(participant_id, timestamp) IN (?)", subQuery).
		Find(&locations)

	if result.Error != nil {
		return nil, result.Error
	}

	return locations, nil
}

func (r *locationRepository) GetHistory(ctx context.Context, participantID uuid.UUID, entityID uuid.UUID, from, to time.Time) ([]*domain.Location, error) {
	var locations []*domain.Location

	result := r.db.WithContext(ctx).
		Where("participant_id = ? AND entity_id = ?", participantID, entityID).
		Where("timestamp >= ? AND timestamp <= ?", from, to).
		Order("timestamp ASC").
		Find(&locations)

	if result.Error != nil {
		return nil, result.Error
	}

	return locations, nil
}
