package service

import (
	"context"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SchedulerService define os métodos do serviço de agendamento
type SchedulerService interface {
	// Criar agendamento
	Create(ctx context.Context, input *domain.CreateSchedulerInput, orgID uuid.UUID) (*domain.Scheduler, error)

	// Buscar por ID
	GetByID(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (*domain.Scheduler, error)

	// Cancelar agendamento
	Cancel(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error

	// Processar tasks pendentes (chamado pelo worker)
	ProcessPendingTasks(ctx context.Context, limit int) (int, error)
}

type schedulerServiceImpl struct {
	schedulerRepo       repository.SchedulerRepository
	participantRepo     repository.ParticipantRepository
	eventRepo           repository.EventRepository
	notificationService NotificationService
	logger              *zap.Logger
}

func NewSchedulerService(
	schedulerRepo repository.SchedulerRepository,
	participantRepo repository.ParticipantRepository,
	eventRepo repository.EventRepository,
	notificationService NotificationService,
	logger *zap.Logger,
) SchedulerService {
	return &schedulerServiceImpl{
		schedulerRepo:       schedulerRepo,
		participantRepo:     participantRepo,
		eventRepo:           eventRepo,
		notificationService: notificationService,
		logger:              logger,
	}
}

// Create cria um novo agendamento
func (s *schedulerServiceImpl) Create(ctx context.Context, input *domain.CreateSchedulerInput, orgID uuid.UUID) (*domain.Scheduler, error) {
	scheduler := &domain.Scheduler{
		ID:          uuid.New(),
		EntityID:    orgID,
		EventID:     input.EventID,
		InstanceID:  input.InstanceID,
		Action:      input.Action,
		Status:      domain.SchedulerStatusPending,
		ScheduledAt: input.ScheduledAt,
		Retries:     0,
		MaxRetries:  input.MaxRetries,
		Metadata:    input.Metadata,
	}

	if scheduler.MaxRetries == 0 {
		scheduler.MaxRetries = 3 // Default
	}

	if err := s.schedulerRepo.Create(ctx, scheduler); err != nil {
		return nil, err
	}

	s.logger.Info("Scheduler created",
		zap.String("id", scheduler.ID.String()),
		zap.String("action", string(scheduler.Action)),
		zap.Time("scheduled_at", scheduler.ScheduledAt),
	)

	return scheduler, nil
}

// GetByID busca um agendamento por ID
func (s *schedulerServiceImpl) GetByID(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (*domain.Scheduler, error) {
	return s.schedulerRepo.GetByID(ctx, id, orgID)
}

// Cancel cancela um agendamento pendente
func (s *schedulerServiceImpl) Cancel(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error {
	scheduler, err := s.schedulerRepo.GetByID(ctx, id, orgID)
	if err != nil {
		return err
	}

	if scheduler.Status != domain.SchedulerStatusPending {
		return domain.ErrInvalidInput
	}

	scheduler.Status = domain.SchedulerStatusSkipped
	return s.schedulerRepo.Update(ctx, scheduler)
}

// ProcessPendingTasks processa as tasks pendentes
func (s *schedulerServiceImpl) ProcessPendingTasks(ctx context.Context, limit int) (int, error) {
	// Buscar tasks pendentes que já passaram do horário
	tasks, err := s.schedulerRepo.ListPending(ctx, time.Now(), limit)
	if err != nil {
		return 0, err
	}

	if len(tasks) == 0 {
		return 0, nil
	}

	s.logger.Debug("Found pending tasks", zap.Int("count", len(tasks)))

	processed := 0
	for _, task := range tasks {
		if err := s.processTask(ctx, task); err != nil {
			s.logger.Error("Failed to process task",
				zap.String("task_id", task.ID.String()),
				zap.String("action", string(task.Action)),
				zap.Error(err),
			)

			// Incrementar retries
			_ = s.schedulerRepo.IncrementRetries(ctx, task.ID, task.EntityID)

			// Se excedeu max retries, marcar como falha
			if task.Retries+1 >= task.MaxRetries {
				_ = s.schedulerRepo.MarkAsFailed(ctx, task.ID, task.EntityID, err.Error())
			}
			continue
		}

		// Marcar como processado
		if err := s.schedulerRepo.MarkAsProcessed(ctx, task.ID, task.EntityID); err != nil {
			s.logger.Error("Failed to mark task as processed",
				zap.String("task_id", task.ID.String()),
				zap.Error(err),
			)
		}

		processed++
	}

	return processed, nil
}

// processTask processa uma task individual
func (s *schedulerServiceImpl) processTask(ctx context.Context, task *domain.Scheduler) error {
	s.logger.Info("Processing task",
		zap.String("task_id", task.ID.String()),
		zap.String("action", string(task.Action)),
		zap.String("event_id", task.EventID.String()),
	)

	switch task.Action {
	case domain.SchedulerActionConfirmation:
		return s.processConfirmation(ctx, task)

	case domain.SchedulerActionReminder:
		return s.processReminder(ctx, task)

	case domain.SchedulerActionClosure:
		return s.processClosure(ctx, task)

	case domain.SchedulerActionLocation:
		return s.processLocationRequest(ctx, task)

	default:
		s.logger.Warn("Unknown scheduler action", zap.String("action", string(task.Action)))
		return nil
	}
}

// processConfirmation envia pedido de confirmação para participantes
func (s *schedulerServiceImpl) processConfirmation(ctx context.Context, task *domain.Scheduler) error {
	// Buscar evento
	event, err := s.eventRepo.GetByID(ctx, task.EventID, task.EntityID)
	if err != nil {
		return err
	}

	// Buscar participantes pendentes
	participants, _, err := s.participantRepo.ListByEvent(ctx, task.EventID, task.EntityID, 1, 1000)
	if err != nil {
		return err
	}

	// Filtrar apenas pendentes
	for _, p := range participants {
		if p.Status != domain.ParticipantStatusPending {
			continue
		}

		if err := s.notificationService.SendConfirmationRequest(ctx, event, p); err != nil {
			s.logger.Error("Failed to send confirmation",
				zap.String("participant_id", p.ID.String()),
				zap.Error(err),
			)
			// Continua com os outros participantes
		}
	}

	return nil
}

// processReminder envia lembretes para participantes confirmados
func (s *schedulerServiceImpl) processReminder(ctx context.Context, task *domain.Scheduler) error {
	// Buscar evento
	event, err := s.eventRepo.GetByID(ctx, task.EventID, task.EntityID)
	if err != nil {
		return err
	}

	// Buscar participantes confirmados
	participants, _, err := s.participantRepo.ListByEvent(ctx, task.EventID, task.EntityID, 1, 1000)
	if err != nil {
		return err
	}

	// Filtrar apenas confirmados
	for _, p := range participants {
		if p.Status != domain.ParticipantStatusConfirmed {
			continue
		}

		if err := s.notificationService.SendReminder(ctx, event, p); err != nil {
			s.logger.Error("Failed to send reminder",
				zap.String("participant_id", p.ID.String()),
				zap.Error(err),
			)
		}
	}

	return nil
}

// processClosure fecha o evento
func (s *schedulerServiceImpl) processClosure(ctx context.Context, task *domain.Scheduler) error {
	// Atualizar status do evento para completed
	return s.eventRepo.Update(ctx, task.EventID, task.EntityID, &domain.UpdateEventInput{
		Status: func() *domain.EventStatus { s := domain.EventStatusCompleted; return &s }(),
	})
}

// processLocationRequest solicita localização dos participantes
func (s *schedulerServiceImpl) processLocationRequest(ctx context.Context, task *domain.Scheduler) error {
	// Buscar evento
	event, err := s.eventRepo.GetByID(ctx, task.EventID, task.EntityID)
	if err != nil {
		return err
	}

	// Buscar participantes confirmados que ainda não fizeram check-in
	participants, _, err := s.participantRepo.ListByEvent(ctx, task.EventID, task.EntityID, 1, 1000)
	if err != nil {
		return err
	}

	for _, p := range participants {
		if p.Status != domain.ParticipantStatusConfirmed {
			continue
		}

		if err := s.notificationService.SendLocationRequest(ctx, event, p); err != nil {
			s.logger.Error("Failed to send location request",
				zap.String("participant_id", p.ID.String()),
				zap.Error(err),
			)
		}
	}

	return nil
}
