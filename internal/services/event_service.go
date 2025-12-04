package services

import (
	"context"

	"event-coming/internal/models"
	"event-coming/internal/repositories"
	"github.com/google/uuid"
)

type EventService struct {
	Events       repositories.EventRepository
	Consolidated repositories.ConsolidatedRepository
}

func NewEventService(er *repositories.EventRepository, cr *repositories.ConsolidatedRepository) *EventService {
	return &EventService{Events: *er, Consolidated: *cr}
}

func (s *EventService) CreateEvent(ctx context.Context, e *models.Event) error {
	return s.Events.Create(ctx, e)
}

func (s *EventService) UpdateEventStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.Events.UpdateStatus(ctx, id, status)
}

func (s *EventService) GetPendingEvents(ctx context.Context, limit, offset int) ([]models.Event, error) {
	filters := map[string]interface{}{"status": "pending"}
	return s.Events.List(ctx, filters, limit, offset)
}

func (s *EventService) FinalizeConsolidated(ctx context.Context, c *models.Consolidated) error {
	return s.Consolidated.Create(ctx, c)
}
