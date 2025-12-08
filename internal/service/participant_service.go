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

// ParticipantService gerencia operações de participantes
type ParticipantService struct {
	participantRepo repository.ParticipantRepository
	eventRepo       repository.EventRepository
}

// NewParticipantService cria um novo serviço de participantes
func NewParticipantService(
	participantRepo repository.ParticipantRepository,
	eventRepo repository.EventRepository,
) *ParticipantService {
	return &ParticipantService{
		participantRepo: participantRepo,
		eventRepo:       eventRepo,
	}
}

// Create cria um novo participante vinculado a um evento
func (s *ParticipantService) Create(ctx context.Context, orgID, eventID uuid.UUID, req *dto.CreateParticipantRequest) (*dto.ParticipantResponse, error) {
	// Verificar se o evento existe
	event, err := s.eventRepo.GetByID(ctx, eventID, orgID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Verificar se já existe participante com mesmo telefone neste evento
	existing, err := s.participantRepo.GetByPhoneNumber(ctx, req.PhoneNumber, eventID, orgID)
	if err != nil && err != domain.ErrNotFound {
		return nil, fmt.Errorf("failed to check existing participant: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("participant with this phone number already exists in this event")
	}

	// Criar participante
	participant := &domain.Participant{
		ID:             uuid.New(),
		EventID:        event.ID,
		InstanceID:     req.InstanceID,
		OrganizationID: orgID,
		Name:           req.Name,
		PhoneNumber:    req.PhoneNumber,
		Email:          req.Email,
		Status:         domain.ParticipantStatusPending,
		Metadata:       req.Metadata,
	}

	if err := s.participantRepo.Create(ctx, participant); err != nil {
		return nil, fmt.Errorf("failed to create participant: %w", err)
	}

	return dto.ToParticipantResponse(participant), nil
}

// GetByID busca um participante por ID
func (s *ParticipantService) GetByID(ctx context.Context, orgID, participantID uuid.UUID) (*dto.ParticipantResponse, error) {
	participant, err := s.participantRepo.GetByID(ctx, participantID, orgID)
	if err != nil {
		return nil, err
	}
	return dto.ToParticipantResponse(participant), nil
}

// Update atualiza um participante
func (s *ParticipantService) Update(ctx context.Context, orgID, participantID uuid.UUID, req *dto.UpdateParticipantRequest) (*dto.ParticipantResponse, error) {
	// Verificar se existe
	participant, err := s.participantRepo.GetByID(ctx, participantID, orgID)
	if err != nil {
		return nil, err
	}

	// Preparar input de atualização
	input := &domain.UpdateParticipantInput{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Status:      req.Status,
		Metadata:    req.Metadata,
	}

	// Atualizar timestamps de status
	if req.Status != nil {
		now := time.Now()
		switch *req.Status {
		case domain.ParticipantStatusConfirmed:
			if participant.ConfirmedAt == nil {
				participant.ConfirmedAt = &now
			}
		case domain.ParticipantStatusCheckedIn:
			if participant.CheckedInAt == nil {
				participant.CheckedInAt = &now
			}
		}
	}

	if err := s.participantRepo.Update(ctx, participantID, orgID, input); err != nil {
		return nil, fmt.Errorf("failed to update participant: %w", err)
	}

	// Buscar participante atualizado
	updated, err := s.participantRepo.GetByID(ctx, participantID, orgID)
	if err != nil {
		return nil, err
	}

	return dto.ToParticipantResponse(updated), nil
}

// Delete remove um participante
func (s *ParticipantService) Delete(ctx context.Context, orgID, participantID uuid.UUID) error {
	return s.participantRepo.Delete(ctx, participantID, orgID)
}

// ListByEvent lista participantes de um evento
func (s *ParticipantService) ListByEvent(ctx context.Context, orgID, eventID uuid.UUID, page, perPage int) ([]*dto.ParticipantResponse, int64, error) {
	// Verificar se o evento existe
	_, err := s.eventRepo.GetByID(ctx, eventID, orgID)
	if err != nil {
		return nil, 0, err
	}

	participants, total, err := s.participantRepo.ListByEvent(ctx, eventID, orgID, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list participants: %w", err)
	}

	responses := make([]*dto.ParticipantResponse, len(participants))
	for i, p := range participants {
		responses[i] = dto.ToParticipantResponse(p)
	}

	return responses, total, nil
}

// UpdateStatus atualiza apenas o status do participante
func (s *ParticipantService) UpdateStatus(ctx context.Context, orgID, participantID uuid.UUID, status domain.ParticipantStatus) error {
	return s.participantRepo.UpdateStatus(ctx, participantID, orgID, status)
}

// ConfirmParticipant confirma a participação
func (s *ParticipantService) ConfirmParticipant(ctx context.Context, orgID, participantID uuid.UUID) (*dto.ParticipantResponse, error) {
	status := domain.ParticipantStatusConfirmed
	return s.Update(ctx, orgID, participantID, &dto.UpdateParticipantRequest{
		Status: &status,
	})
}

// CheckInParticipant faz check-in do participante
func (s *ParticipantService) CheckInParticipant(ctx context.Context, orgID, participantID uuid.UUID) (*dto.ParticipantResponse, error) {
	status := domain.ParticipantStatusCheckedIn
	return s.Update(ctx, orgID, participantID, &dto.UpdateParticipantRequest{
		Status: &status,
	})
}

// BatchCreate cria múltiplos participantes de uma vez
func (s *ParticipantService) BatchCreate(ctx context.Context, orgID, eventID uuid.UUID, req *dto.BatchCreateParticipantsRequest) ([]*dto.ParticipantResponse, []error) {
	// Verificar se o evento existe
	_, err := s.eventRepo.GetByID(ctx, eventID, orgID)
	if err != nil {
		return nil, []error{fmt.Errorf("event not found: %w", err)}
	}

	var responses []*dto.ParticipantResponse
	var errors []error

	for i, pReq := range req.Participants {
		resp, err := s.Create(ctx, orgID, eventID, &pReq)
		if err != nil {
			errors = append(errors, fmt.Errorf("participant[%d]: %w", i, err))
			continue
		}
		responses = append(responses, resp)
	}

	return responses, errors
}
