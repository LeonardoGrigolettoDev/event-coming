package service

import (
	"context"
	"fmt"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/dto"
	"event-coming/internal/repository"

	"github.com/google/uuid"
)

// EventService gerencia operações de eventos
type EventService struct {
	eventRepo       repository.EventRepository
	schedulerRepo   repository.SchedulerRepository
	participantRepo repository.ParticipantRepository
}

// NewEventService cria um novo serviço de eventos
func NewEventService(
	eventRepo repository.EventRepository,
	schedulerRepo repository.SchedulerRepository,
	participantRepo repository.ParticipantRepository,
) *EventService {
	return &EventService{
		eventRepo:       eventRepo,
		schedulerRepo:   schedulerRepo,
		participantRepo: participantRepo,
	}
}

// Create cria um novo evento com schedulers e participants opcionais
func (s *EventService) Create(ctx context.Context, entID, userID uuid.UUID, req *dto.CreateEventRequest) (*dto.EventResponse, error) {
	// Criar evento
	event := &domain.Event{
		ID:                   uuid.New(),
		EntityID:             entID,
		Name:                 req.Name,
		Description:          req.Description,
		Type:                 req.Type,
		Status:               domain.EventStatusDraft,
		LocationLat:          req.LocationLat,
		LocationLng:          req.LocationLng,
		LocationAddress:      req.LocationAddress,
		StartTime:            req.StartTime,
		EndTime:              req.EndTime,
		RRuleString:          req.RRuleString,
		ConfirmationDeadline: req.ConfirmationDeadline,
		CreatedBy:            userID,
	}

	if err := s.eventRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	response := dto.ToEventResponse(event)

	// Criar schedulers
	schedulersCreated := 0
	if req.Scheduler != nil {
		count, err := s.createSchedulers(ctx, entID, event, req.Scheduler)
		if err != nil {
			fmt.Printf("Warning: failed to create some schedulers: %v\n", err)
		}
		schedulersCreated = count
	} else {
		count, _ := s.createDefaultSchedulers(ctx, entID, event)
		schedulersCreated = count
	}
	response.SchedulersCreated = schedulersCreated

	// Criar participants
	if len(req.Participants) > 0 {
		participants, _ := s.createParticipants(ctx, entID, event.ID, req.Participants)
		response.Participants = participants
	}

	return response, nil
}

// createSchedulers cria schedulers baseado na configuração
func (s *EventService) createSchedulers(ctx context.Context, entID uuid.UUID, event *domain.Event, config *dto.SchedulerConfig) (int, error) {
	var count int
	var lastErr error

	// Scheduler de confirmação
	if config.SendConfirmation {
		scheduledAt := event.StartTime.Add(-24 * time.Hour)
		if config.ConfirmationTime != nil {
			scheduledAt = *config.ConfirmationTime
		}

		scheduler := &domain.Scheduler{
			ID:          uuid.New(),
			EntityID:    entID,
			EventID:     event.ID,
			Action:      domain.SchedulerActionConfirmation,
			Status:      domain.SchedulerStatusPending,
			ScheduledAt: scheduledAt,
			MaxRetries:  3,
			Metadata: map[string]interface{}{
				"event_name": event.Name,
			},
		}

		if err := s.schedulerRepo.Create(ctx, scheduler); err != nil {
			lastErr = err
		} else {
			count++
		}
	}

	// Scheduler de lembrete
	if config.SendReminder {
		scheduledAt := event.StartTime.Add(-2 * time.Hour)
		if config.ReminderTime != nil {
			scheduledAt = *config.ReminderTime
		} else if config.ReminderBeforeHours != nil {
			scheduledAt = event.StartTime.Add(-time.Duration(*config.ReminderBeforeHours) * time.Hour)
		}

		scheduler := &domain.Scheduler{
			ID:          uuid.New(),
			EntityID:    entID,
			EventID:     event.ID,
			Action:      domain.SchedulerActionReminder,
			Status:      domain.SchedulerStatusPending,
			ScheduledAt: scheduledAt,
			MaxRetries:  3,
			Metadata: map[string]interface{}{
				"event_name": event.Name,
			},
		}

		if err := s.schedulerRepo.Create(ctx, scheduler); err != nil {
			lastErr = err
		} else {
			count++
		}
	}

	// Scheduler de rastreamento de localização
	if config.TrackLocation {
		scheduledAt := event.StartTime.Add(-1 * time.Hour)
		if config.LocationTrackingTime != nil {
			scheduledAt = *config.LocationTrackingTime
		}

		scheduler := &domain.Scheduler{
			ID:          uuid.New(),
			EntityID:    entID,
			EventID:     event.ID,
			Action:      domain.SchedulerActionLocation,
			Status:      domain.SchedulerStatusPending,
			ScheduledAt: scheduledAt,
			MaxRetries:  3,
			Metadata: map[string]interface{}{
				"event_name":   event.Name,
				"location_lat": event.LocationLat,
				"location_lng": event.LocationLng,
			},
		}

		if err := s.schedulerRepo.Create(ctx, scheduler); err != nil {
			lastErr = err
		} else {
			count++
		}
	}

	// Scheduler de fechamento (sempre criar)
	closureScheduler := &domain.Scheduler{
		ID:          uuid.New(),
		EntityID:    entID,
		EventID:     event.ID,
		Action:      domain.SchedulerActionClosure,
		Status:      domain.SchedulerStatusPending,
		ScheduledAt: event.StartTime,
		MaxRetries:  3,
		Metadata: map[string]interface{}{
			"event_name": event.Name,
		},
	}
	if event.EndTime != nil {
		closureScheduler.ScheduledAt = *event.EndTime
	}

	if err := s.schedulerRepo.Create(ctx, closureScheduler); err != nil {
		lastErr = err
	} else {
		count++
	}

	return count, lastErr
}

// createDefaultSchedulers cria schedulers padrão para um evento
func (s *EventService) createDefaultSchedulers(ctx context.Context, entID uuid.UUID, event *domain.Event) (int, error) {
	config := &dto.SchedulerConfig{
		SendConfirmation: true,
		SendReminder:     true,
		TrackLocation:    true,
	}
	return s.createSchedulers(ctx, entID, event, config)
}

// createParticipants cria participants para o evento
func (s *EventService) createParticipants(ctx context.Context, entID, eventID uuid.UUID, inputs []dto.ParticipantInput) ([]*dto.ParticipantResponse, error) {
	var participants []*dto.ParticipantResponse
	var lastErr error

	for _, input := range inputs {
		participant := &domain.Participant{
			ID:       uuid.New(),
			EventID:  eventID,
			EntityID: entID,
			Status:   domain.ParticipantStatusPending,
			Metadata: input.Metadata,
		}

		if err := s.participantRepo.Create(ctx, participant); err != nil {
			lastErr = err
			continue
		}

		participants = append(participants, dto.ToParticipantResponse(participant))
	}

	return participants, lastErr
}

// GetByID busca um evento por ID
func (s *EventService) GetByID(ctx context.Context, entID, eventID uuid.UUID) (*dto.EventResponse, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID, entID)
	if err != nil {
		return nil, err
	}
	return dto.ToEventResponse(event), nil
}

// GetByIDWithParticipants busca um evento com seus participants
func (s *EventService) GetByIDWithParticipants(ctx context.Context, entID, eventID uuid.UUID) (*dto.EventResponse, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID, entID)
	if err != nil {
		return nil, err
	}

	response := dto.ToEventResponse(event)

	// Buscar participants
	participants, _, err := s.participantRepo.ListByEvent(ctx, eventID, entID, 1, 1000)
	if err == nil {
		for _, p := range participants {
			response.Participants = append(response.Participants, dto.ToParticipantResponse(p))
		}
	}

	return response, nil
}

// Update atualiza um evento
func (s *EventService) Update(ctx context.Context, entID, eventID uuid.UUID, req *dto.UpdateEventRequest) (*dto.EventResponse, error) {
	_, err := s.eventRepo.GetByID(ctx, eventID, entID)
	if err != nil {
		return nil, err
	}

	input := &domain.UpdateEventInput{
		Name:                 req.Name,
		Description:          req.Description,
		Status:               req.Status,
		LocationLat:          req.LocationLat,
		LocationLng:          req.LocationLng,
		LocationAddress:      req.LocationAddress,
		StartTime:            req.StartTime,
		EndTime:              req.EndTime,
		ConfirmationDeadline: req.ConfirmationDeadline,
	}

	if err := s.eventRepo.Update(ctx, eventID, entID, input); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	updated, err := s.eventRepo.GetByID(ctx, eventID, entID)
	if err != nil {
		return nil, err
	}

	return dto.ToEventResponse(updated), nil
}

// Delete remove um evento
func (s *EventService) Delete(ctx context.Context, entID, eventID uuid.UUID) error {
	return s.eventRepo.Delete(ctx, eventID, entID)
}

// List lista eventos de uma organização
func (s *EventService) List(ctx context.Context, entID uuid.UUID, page, perPage int) ([]*dto.EventResponse, int64, error) {
	events, total, err := s.eventRepo.List(ctx, entID, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list events: %w", err)
	}

	responses := make([]*dto.EventResponse, len(events))
	for i, e := range events {
		responses[i] = dto.ToEventResponse(e)
	}

	return responses, total, nil
}

// ListByStatus lista eventos por status
func (s *EventService) ListByStatus(ctx context.Context, entID uuid.UUID, status domain.EventStatus, page, perPage int) ([]*dto.EventResponse, int64, error) {
	events, total, err := s.eventRepo.ListByStatus(ctx, entID, status, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list events: %w", err)
	}

	responses := make([]*dto.EventResponse, len(events))
	for i, e := range events {
		responses[i] = dto.ToEventResponse(e)
	}

	return responses, total, nil
}

// Activate ativa um evento
func (s *EventService) Activate(ctx context.Context, entID, eventID uuid.UUID) (*dto.EventResponse, error) {
	status := domain.EventStatusActive
	return s.Update(ctx, entID, eventID, &dto.UpdateEventRequest{Status: &status})
}

// Cancel cancela um evento
func (s *EventService) Cancel(ctx context.Context, entID, eventID uuid.UUID) (*dto.EventResponse, error) {
	status := domain.EventStatusCancelled
	return s.Update(ctx, entID, eventID, &dto.UpdateEventRequest{Status: &status})
}

// Complete marca um evento como completo
func (s *EventService) Complete(ctx context.Context, entID, eventID uuid.UUID) (*dto.EventResponse, error) {
	status := domain.EventStatusCompleted
	return s.Update(ctx, entID, eventID, &dto.UpdateEventRequest{Status: &status})
}
