package mocks

import (
	"context"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockEventRepository is a mock implementation of EventRepository
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Create(ctx context.Context, event *domain.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) GetByID(ctx context.Context, id uuid.UUID, entityID uuid.UUID) (*domain.Event, error) {
	args := m.Called(ctx, id, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Event), args.Error(1)
}

func (m *MockEventRepository) Update(ctx context.Context, id uuid.UUID, entityID uuid.UUID, input *domain.UpdateEventInput) error {
	args := m.Called(ctx, id, entityID, input)
	return args.Error(0)
}

func (m *MockEventRepository) Delete(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error {
	args := m.Called(ctx, id, entityID)
	return args.Error(0)
}

func (m *MockEventRepository) List(ctx context.Context, entityID uuid.UUID, page, perPage int) ([]*domain.Event, int64, error) {
	args := m.Called(ctx, entityID, page, perPage)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Event), args.Get(1).(int64), args.Error(2)
}

func (m *MockEventRepository) ListByStatus(ctx context.Context, entityID uuid.UUID, status domain.EventStatus, page, perPage int) ([]*domain.Event, int64, error) {
	args := m.Called(ctx, entityID, status, page, perPage)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Event), args.Get(1).(int64), args.Error(2)
}

func (m *MockEventRepository) CreateInstance(ctx context.Context, instance *domain.EventInstance) error {
	args := m.Called(ctx, instance)
	return args.Error(0)
}

func (m *MockEventRepository) GetInstanceByID(ctx context.Context, id uuid.UUID, entityID uuid.UUID) (*domain.EventInstance, error) {
	args := m.Called(ctx, id, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EventInstance), args.Error(1)
}

func (m *MockEventRepository) ListInstances(ctx context.Context, eventID uuid.UUID, entityID uuid.UUID) ([]*domain.EventInstance, error) {
	args := m.Called(ctx, eventID, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.EventInstance), args.Error(1)
}

// MockParticipantRepository is a mock implementation of ParticipantRepository
type MockParticipantRepository struct {
	mock.Mock
}

func (m *MockParticipantRepository) Create(ctx context.Context, participant *domain.Participant) error {
	args := m.Called(ctx, participant)
	return args.Error(0)
}

func (m *MockParticipantRepository) GetByID(ctx context.Context, id uuid.UUID, entityID uuid.UUID) (*domain.Participant, error) {
	args := m.Called(ctx, id, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Participant), args.Error(1)
}

func (m *MockParticipantRepository) Update(ctx context.Context, id uuid.UUID, entityID uuid.UUID, input *domain.UpdateParticipantInput) error {
	args := m.Called(ctx, id, entityID, input)
	return args.Error(0)
}

func (m *MockParticipantRepository) Delete(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error {
	args := m.Called(ctx, id, entityID)
	return args.Error(0)
}

func (m *MockParticipantRepository) ListByEvent(ctx context.Context, eventID uuid.UUID, entityID uuid.UUID, page, perPage int) ([]*domain.Participant, int64, error) {
	args := m.Called(ctx, eventID, entityID, page, perPage)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Participant), args.Get(1).(int64), args.Error(2)
}

func (m *MockParticipantRepository) ListByEventInstance(ctx context.Context, instanceID uuid.UUID, entityID uuid.UUID, page, perPage int) ([]*domain.Participant, int64, error) {
	args := m.Called(ctx, instanceID, entityID, page, perPage)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Participant), args.Get(1).(int64), args.Error(2)
}

func (m *MockParticipantRepository) UpdateStatus(ctx context.Context, id uuid.UUID, entityID uuid.UUID, status domain.ParticipantStatus) error {
	args := m.Called(ctx, id, entityID, status)
	return args.Error(0)
}

func (m *MockParticipantRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string, eventID uuid.UUID, entityID uuid.UUID) (*domain.Participant, error) {
	args := m.Called(ctx, phoneNumber, eventID, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Participant), args.Error(1)
}

func (m *MockParticipantRepository) GetActiveByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Participant, error) {
	args := m.Called(ctx, phoneNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Participant), args.Error(1)
}

// MockLocationRepository is a mock implementation of LocationRepository
type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) Create(ctx context.Context, location *domain.Location) error {
	args := m.Called(ctx, location)
	return args.Error(0)
}

func (m *MockLocationRepository) BatchCreate(ctx context.Context, locations []*domain.Location) error {
	args := m.Called(ctx, locations)
	return args.Error(0)
}

func (m *MockLocationRepository) GetLatestByParticipant(ctx context.Context, participantID uuid.UUID, entityID uuid.UUID) (*domain.Location, error) {
	args := m.Called(ctx, participantID, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Location), args.Error(1)
}

func (m *MockLocationRepository) GetLatestByEvent(ctx context.Context, eventID uuid.UUID, entityID uuid.UUID) ([]*domain.Location, error) {
	args := m.Called(ctx, eventID, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Location), args.Error(1)
}

func (m *MockLocationRepository) GetHistory(ctx context.Context, participantID uuid.UUID, entityID uuid.UUID, from, to time.Time) ([]*domain.Location, error) {
	args := m.Called(ctx, participantID, entityID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Location), args.Error(1)
}

// MockSchedulerRepository is a mock implementation of SchedulerRepository
type MockSchedulerRepository struct {
	mock.Mock
}

func (m *MockSchedulerRepository) Create(ctx context.Context, scheduler *domain.Scheduler) error {
	args := m.Called(ctx, scheduler)
	return args.Error(0)
}

func (m *MockSchedulerRepository) GetByID(ctx context.Context, id uuid.UUID, entityID uuid.UUID) (*domain.Scheduler, error) {
	args := m.Called(ctx, id, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Scheduler), args.Error(1)
}

func (m *MockSchedulerRepository) Update(ctx context.Context, scheduler *domain.Scheduler) error {
	args := m.Called(ctx, scheduler)
	return args.Error(0)
}

func (m *MockSchedulerRepository) Delete(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error {
	args := m.Called(ctx, id, entityID)
	return args.Error(0)
}

func (m *MockSchedulerRepository) ListPending(ctx context.Context, before time.Time, limit int) ([]*domain.Scheduler, error) {
	args := m.Called(ctx, before, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Scheduler), args.Error(1)
}

func (m *MockSchedulerRepository) MarkAsProcessed(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error {
	args := m.Called(ctx, id, entityID)
	return args.Error(0)
}

func (m *MockSchedulerRepository) MarkAsFailed(ctx context.Context, id uuid.UUID, entityID uuid.UUID, errorMsg string) error {
	args := m.Called(ctx, id, entityID, errorMsg)
	return args.Error(0)
}

func (m *MockSchedulerRepository) IncrementRetries(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error {
	args := m.Called(ctx, id, entityID)
	return args.Error(0)
}
