package mocks

import (
	"context"

	"event-coming/internal/domain"
	"event-coming/internal/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegisterResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *MockAuthService) Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RefreshResponse), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, req dto.LogoutRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuthService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ForgotPasswordResponse), args.Error(1)
}

func (m *MockAuthService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ResetPasswordResponse), args.Error(1)
}

// MockEntityService is a mock implementation of EntityService
type MockEntityService struct {
	mock.Mock
}

func (m *MockEntityService) Create(ctx context.Context, req *dto.CreateEntityRequest) (*dto.EntityResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EntityResponse), args.Error(1)
}

func (m *MockEntityService) GetByID(ctx context.Context, id uuid.UUID) (*dto.EntityResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EntityResponse), args.Error(1)
}

func (m *MockEntityService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateEntityRequest) (*dto.EntityResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EntityResponse), args.Error(1)
}

func (m *MockEntityService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEntityService) List(ctx context.Context, page, perPage int) ([]*dto.EntityResponse, int64, error) {
	args := m.Called(ctx, page, perPage)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*dto.EntityResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockEntityService) ListByParent(ctx context.Context, parentID uuid.UUID, page, perPage int) ([]*dto.EntityResponse, int64, error) {
	args := m.Called(ctx, parentID, page, perPage)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*dto.EntityResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockEntityService) GetByDocument(ctx context.Context, document string) (*dto.EntityResponse, error) {
	args := m.Called(ctx, document)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EntityResponse), args.Error(1)
}

// MockEventService is a mock implementation of EventService
type MockEventService struct {
	mock.Mock
}

func (m *MockEventService) Create(ctx context.Context, req *dto.CreateEventRequest) (*dto.EventResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EventResponse), args.Error(1)
}

func (m *MockEventService) GetByID(ctx context.Context, id uuid.UUID) (*dto.EventResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EventResponse), args.Error(1)
}

func (m *MockEventService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateEventRequest) (*dto.EventResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EventResponse), args.Error(1)
}

func (m *MockEventService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventService) List(ctx context.Context, entityID uuid.UUID, page, perPage int) ([]*dto.EventResponse, int64, error) {
	args := m.Called(ctx, entityID, page, perPage)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*dto.EventResponse), args.Get(1).(int64), args.Error(2)
}

// MockParticipantService is a mock implementation of ParticipantService
type MockParticipantService struct {
	mock.Mock
}

func (m *MockParticipantService) Create(ctx context.Context, req *dto.CreateParticipantRequest) (*dto.ParticipantResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ParticipantResponse), args.Error(1)
}

func (m *MockParticipantService) GetByID(ctx context.Context, id uuid.UUID) (*dto.ParticipantResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ParticipantResponse), args.Error(1)
}

func (m *MockParticipantService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateParticipantRequest) (*dto.ParticipantResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ParticipantResponse), args.Error(1)
}

func (m *MockParticipantService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockParticipantService) ListByEvent(ctx context.Context, eventID uuid.UUID, page, perPage int) ([]*dto.ParticipantResponse, int64, error) {
	args := m.Called(ctx, eventID, page, perPage)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*dto.ParticipantResponse), args.Get(1).(int64), args.Error(2)
}

// MockLocationService is a mock implementation of LocationService
type MockLocationService struct {
	mock.Mock
}

func (m *MockLocationService) Create(ctx context.Context, req *dto.CreateLocationRequest) (*dto.LocationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LocationResponse), args.Error(1)
}

func (m *MockLocationService) GetLatestByParticipant(ctx context.Context, participantID uuid.UUID) (*dto.LocationResponse, error) {
	args := m.Called(ctx, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LocationResponse), args.Error(1)
}

func (m *MockLocationService) GetLatestByEvent(ctx context.Context, eventID uuid.UUID) ([]*dto.LocationResponse, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.LocationResponse), args.Error(1)
}

func (m *MockLocationService) CalculateETA(ctx context.Context, participantID uuid.UUID) (*dto.ETAResponse, error) {
	args := m.Called(ctx, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ETAResponse), args.Error(1)
}

// MockNotificationService is a mock implementation of NotificationService
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendEventInvitation(ctx context.Context, participant *domain.Participant, event *domain.Event) error {
	args := m.Called(ctx, participant, event)
	return args.Error(0)
}

func (m *MockNotificationService) SendEventReminder(ctx context.Context, participant *domain.Participant, event *domain.Event) error {
	args := m.Called(ctx, participant, event)
	return args.Error(0)
}

func (m *MockNotificationService) SendLocationRequest(ctx context.Context, participant *domain.Participant, event *domain.Event) error {
	args := m.Called(ctx, participant, event)
	return args.Error(0)
}

func (m *MockNotificationService) SendEventUpdate(ctx context.Context, participant *domain.Participant, event *domain.Event, updateType string) error {
	args := m.Called(ctx, participant, event, updateType)
	return args.Error(0)
}

// MockSchedulerService is a mock implementation of SchedulerService
type MockSchedulerService struct {
	mock.Mock
}

func (m *MockSchedulerService) ScheduleEventNotifications(ctx context.Context, event *domain.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockSchedulerService) CancelEventNotifications(ctx context.Context, eventID uuid.UUID, entityID uuid.UUID) error {
	args := m.Called(ctx, eventID, entityID)
	return args.Error(0)
}

func (m *MockSchedulerService) ProcessPendingTasks(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockETAService is a mock implementation of ETAService
type MockETAService struct {
	mock.Mock
}

func (m *MockETAService) CalculateETA(ctx context.Context, locations []*domain.Location, targetLat, targetLng float64) (*dto.ETAResponse, error) {
	args := m.Called(ctx, locations, targetLat, targetLng)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ETAResponse), args.Error(1)
}
