package services

import (
	"context"

	"github.com/LeonardoGrigolettoDev/event-coming.git/internal/models"
	"github.com/LeonardoGrigolettoDev/event-coming.git/internal/repositories"
	"github.com/google/uuid"
)

type SchedulerService struct {
	SchedulerRepo         *repositories.SchedulerRepository
	SchedulerContactsRepo *repositories.SchedulerContactsRepository
	ContactRepo           *repositories.ContactRepository
	EventRepo             *repositories.EventRepository
}

func NewSchedulerService(sr *repositories.SchedulerRepository, scr *repositories.SchedulerContactsRepository, cr *repositories.ContactRepository, er *repositories.EventRepository) *SchedulerService {
	return &SchedulerService{SchedulerRepo: sr, SchedulerContactsRepo: scr, ContactRepo: cr, EventRepo: er}
}

func (s *SchedulerService) GetSchedulerContacts(ctx context.Context, schedulerID uuid.UUID) ([]models.Contact, error) {
	filters := map[string]interface{}{"scheduler_id": schedulerID}
	scs, err := s.SchedulerContactsRepo.List(ctx, filters, 1000, 0)
	if err != nil {
		return nil, err
	}
	out := make([]models.Contact, 0, len(scs))
	for _, sc := range scs {
		c, err := s.ContactRepo.GetByID(ctx, sc.ContactID)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, nil
}

func (s *SchedulerService) GenerateEventsForScheduler(ctx context.Context, schedulerID uuid.UUID) error {
	contacts, err := s.GetSchedulerContacts(ctx, schedulerID)
	if err != nil {
		return err
	}
	for _, c := range contacts {
		ev := &models.Event{SchedulerID: schedulerID, ContactID: c.ID, EventType: "notification", Status: "pending"}
		if err := s.EventRepo.Create(ctx, ev); err != nil {
			return err
		}
	}
	return nil
}
