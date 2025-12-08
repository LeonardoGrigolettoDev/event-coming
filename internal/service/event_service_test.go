package service

import (
	"context"
	"testing"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/dto"
	"event-coming/internal/testutil/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEventService_Create(t *testing.T) {
	entityID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name      string
		req       *dto.CreateEventRequest
		setupMock func(*mocks.MockEventRepository, *mocks.MockSchedulerRepository, *mocks.MockParticipantRepository)
		wantErr   bool
	}{
		{
			name: "successful creation with default schedulers",
			req: &dto.CreateEventRequest{
				Name:      "Test Event",
				Type:      domain.EventTypeDemand,
				StartTime: time.Now().Add(24 * time.Hour),
			},
			setupMock: func(eventRepo *mocks.MockEventRepository, schedRepo *mocks.MockSchedulerRepository, partRepo *mocks.MockParticipantRepository) {
				eventRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Event")).Return(nil)
				// Default schedulers: confirmation, reminder, location, closure = 4 schedulers
				schedRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Scheduler")).Return(nil).Times(4)
			},
			wantErr: false,
		},
		{
			name: "creation with custom scheduler config",
			req: &dto.CreateEventRequest{
				Name:      "Test Event",
				Type:      domain.EventTypeDemand,
				StartTime: time.Now().Add(24 * time.Hour),
				Scheduler: &dto.SchedulerConfig{
					SendConfirmation: true,
					SendReminder:     false,
					TrackLocation:    false,
				},
			},
			setupMock: func(eventRepo *mocks.MockEventRepository, schedRepo *mocks.MockSchedulerRepository, partRepo *mocks.MockParticipantRepository) {
				eventRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Event")).Return(nil)
				// confirmation + closure = 2 schedulers
				schedRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Scheduler")).Return(nil).Times(2)
			},
			wantErr: false,
		},
		{
			name: "creation with participants",
			req: &dto.CreateEventRequest{
				Name:      "Test Event",
				Type:      domain.EventTypeDemand,
				StartTime: time.Now().Add(24 * time.Hour),
				Participants: []dto.ParticipantInput{
					{Name: "John", PhoneNumber: "+5511999999999"},
					{Name: "Jane", PhoneNumber: "+5511888888888"},
				},
			},
			setupMock: func(eventRepo *mocks.MockEventRepository, schedRepo *mocks.MockSchedulerRepository, partRepo *mocks.MockParticipantRepository) {
				eventRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Event")).Return(nil)
				schedRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Scheduler")).Return(nil).Times(4)
				partRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Participant")).Return(nil).Times(2)
			},
			wantErr: false,
		},
		{
			name: "event creation fails",
			req: &dto.CreateEventRequest{
				Name:      "Test Event",
				Type:      domain.EventTypeDemand,
				StartTime: time.Now().Add(24 * time.Hour),
			},
			setupMock: func(eventRepo *mocks.MockEventRepository, schedRepo *mocks.MockSchedulerRepository, partRepo *mocks.MockParticipantRepository) {
				eventRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Event")).Return(domain.ErrInvalidInput)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventRepo := new(mocks.MockEventRepository)
			schedRepo := new(mocks.MockSchedulerRepository)
			partRepo := new(mocks.MockParticipantRepository)

			tt.setupMock(eventRepo, schedRepo, partRepo)

			svc := NewEventService(eventRepo, schedRepo, partRepo)
			result, err := svc.Create(context.Background(), entityID, userID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Name, result.Name)
			}

			eventRepo.AssertExpectations(t)
		})
	}
}

func TestEventService_GetByID(t *testing.T) {
	entityID := uuid.New()
	eventID := uuid.New()

	tests := []struct {
		name      string
		setupMock func(*mocks.MockEventRepository)
		wantErr   bool
	}{
		{
			name: "event found",
			setupMock: func(eventRepo *mocks.MockEventRepository) {
				event := &domain.Event{
					ID:       eventID,
					EntityID: entityID,
					Name:     "Test Event",
					Type:     domain.EventTypeDemand,
				}
				eventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(event, nil)
			},
			wantErr: false,
		},
		{
			name: "event not found",
			setupMock: func(eventRepo *mocks.MockEventRepository) {
				eventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventRepo := new(mocks.MockEventRepository)
			schedRepo := new(mocks.MockSchedulerRepository)
			partRepo := new(mocks.MockParticipantRepository)

			tt.setupMock(eventRepo)

			svc := NewEventService(eventRepo, schedRepo, partRepo)
			result, err := svc.GetByID(context.Background(), entityID, eventID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			eventRepo.AssertExpectations(t)
		})
	}
}

func TestEventService_Update(t *testing.T) {
	entityID := uuid.New()
	eventID := uuid.New()

	tests := []struct {
		name      string
		req       *dto.UpdateEventRequest
		setupMock func(*mocks.MockEventRepository)
		wantErr   bool
	}{
		{
			name: "successful update",
			req: &dto.UpdateEventRequest{
				Name: eventStrPtr("Updated Event Name"),
			},
			setupMock: func(eventRepo *mocks.MockEventRepository) {
				event := &domain.Event{ID: eventID, EntityID: entityID, Name: "Old Name"}
				eventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(event, nil).Once()
				eventRepo.On("Update", mock.Anything, eventID, entityID, mock.AnythingOfType("*domain.UpdateEventInput")).Return(nil)
				updatedEvent := &domain.Event{ID: eventID, EntityID: entityID, Name: "Updated Event Name"}
				eventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(updatedEvent, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "event not found",
			req: &dto.UpdateEventRequest{
				Name: eventStrPtr("Updated"),
			},
			setupMock: func(eventRepo *mocks.MockEventRepository) {
				eventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventRepo := new(mocks.MockEventRepository)
			schedRepo := new(mocks.MockSchedulerRepository)
			partRepo := new(mocks.MockParticipantRepository)

			tt.setupMock(eventRepo)

			svc := NewEventService(eventRepo, schedRepo, partRepo)
			result, err := svc.Update(context.Background(), entityID, eventID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			eventRepo.AssertExpectations(t)
		})
	}
}

func TestEventService_Delete(t *testing.T) {
	entityID := uuid.New()
	eventID := uuid.New()

	tests := []struct {
		name      string
		setupMock func(*mocks.MockEventRepository)
		wantErr   bool
	}{
		{
			name: "successful delete",
			setupMock: func(eventRepo *mocks.MockEventRepository) {
				eventRepo.On("Delete", mock.Anything, eventID, entityID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "delete fails",
			setupMock: func(eventRepo *mocks.MockEventRepository) {
				eventRepo.On("Delete", mock.Anything, eventID, entityID).Return(domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventRepo := new(mocks.MockEventRepository)
			schedRepo := new(mocks.MockSchedulerRepository)
			partRepo := new(mocks.MockParticipantRepository)

			tt.setupMock(eventRepo)

			svc := NewEventService(eventRepo, schedRepo, partRepo)
			err := svc.Delete(context.Background(), entityID, eventID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			eventRepo.AssertExpectations(t)
		})
	}
}

func TestEventService_List(t *testing.T) {
	entityID := uuid.New()

	tests := []struct {
		name      string
		page      int
		perPage   int
		setupMock func(*mocks.MockEventRepository)
		wantLen   int
		wantTotal int64
		wantErr   bool
	}{
		{
			name:    "successful list",
			page:    1,
			perPage: 10,
			setupMock: func(eventRepo *mocks.MockEventRepository) {
				events := []*domain.Event{
					{ID: uuid.New(), EntityID: entityID, Name: "Event 1"},
					{ID: uuid.New(), EntityID: entityID, Name: "Event 2"},
				}
				eventRepo.On("List", mock.Anything, entityID, 1, 10).Return(events, int64(2), nil)
			},
			wantLen:   2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name:    "empty list",
			page:    1,
			perPage: 10,
			setupMock: func(eventRepo *mocks.MockEventRepository) {
				eventRepo.On("List", mock.Anything, entityID, 1, 10).Return([]*domain.Event{}, int64(0), nil)
			},
			wantLen:   0,
			wantTotal: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventRepo := new(mocks.MockEventRepository)
			schedRepo := new(mocks.MockSchedulerRepository)
			partRepo := new(mocks.MockParticipantRepository)

			tt.setupMock(eventRepo)

			svc := NewEventService(eventRepo, schedRepo, partRepo)
			result, total, err := svc.List(context.Background(), entityID, tt.page, tt.perPage)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
				assert.Equal(t, tt.wantTotal, total)
			}

			eventRepo.AssertExpectations(t)
		})
	}
}

func TestEventService_ListByStatus(t *testing.T) {
	entityID := uuid.New()

	eventRepo := new(mocks.MockEventRepository)
	schedRepo := new(mocks.MockSchedulerRepository)
	partRepo := new(mocks.MockParticipantRepository)

	events := []*domain.Event{
		{ID: uuid.New(), EntityID: entityID, Name: "Active Event 1", Status: domain.EventStatusActive},
	}
	eventRepo.On("ListByStatus", mock.Anything, entityID, domain.EventStatusActive, 1, 10).Return(events, int64(1), nil)

	svc := NewEventService(eventRepo, schedRepo, partRepo)
	result, total, err := svc.ListByStatus(context.Background(), entityID, domain.EventStatusActive, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)

	eventRepo.AssertExpectations(t)
}

func TestEventService_Activate(t *testing.T) {
	entityID := uuid.New()
	eventID := uuid.New()

	eventRepo := new(mocks.MockEventRepository)
	schedRepo := new(mocks.MockSchedulerRepository)
	partRepo := new(mocks.MockParticipantRepository)

	event := &domain.Event{ID: eventID, EntityID: entityID, Name: "Event", Status: domain.EventStatusDraft}
	eventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(event, nil).Once()
	eventRepo.On("Update", mock.Anything, eventID, entityID, mock.AnythingOfType("*domain.UpdateEventInput")).Return(nil)
	activeEvent := &domain.Event{ID: eventID, EntityID: entityID, Name: "Event", Status: domain.EventStatusActive}
	eventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(activeEvent, nil).Once()

	svc := NewEventService(eventRepo, schedRepo, partRepo)
	result, err := svc.Activate(context.Background(), entityID, eventID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.EventStatusActive, result.Status)

	eventRepo.AssertExpectations(t)
}

func TestEventService_Cancel(t *testing.T) {
	entityID := uuid.New()
	eventID := uuid.New()

	eventRepo := new(mocks.MockEventRepository)
	schedRepo := new(mocks.MockSchedulerRepository)
	partRepo := new(mocks.MockParticipantRepository)

	event := &domain.Event{ID: eventID, EntityID: entityID, Name: "Event", Status: domain.EventStatusActive}
	eventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(event, nil).Once()
	eventRepo.On("Update", mock.Anything, eventID, entityID, mock.AnythingOfType("*domain.UpdateEventInput")).Return(nil)
	cancelledEvent := &domain.Event{ID: eventID, EntityID: entityID, Name: "Event", Status: domain.EventStatusCancelled}
	eventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(cancelledEvent, nil).Once()

	svc := NewEventService(eventRepo, schedRepo, partRepo)
	result, err := svc.Cancel(context.Background(), entityID, eventID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.EventStatusCancelled, result.Status)

	eventRepo.AssertExpectations(t)
}

func TestEventService_Complete(t *testing.T) {
	entityID := uuid.New()
	eventID := uuid.New()

	eventRepo := new(mocks.MockEventRepository)
	schedRepo := new(mocks.MockSchedulerRepository)
	partRepo := new(mocks.MockParticipantRepository)

	event := &domain.Event{ID: eventID, EntityID: entityID, Name: "Event", Status: domain.EventStatusActive}
	eventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(event, nil).Once()
	eventRepo.On("Update", mock.Anything, eventID, entityID, mock.AnythingOfType("*domain.UpdateEventInput")).Return(nil)
	completedEvent := &domain.Event{ID: eventID, EntityID: entityID, Name: "Event", Status: domain.EventStatusCompleted}
	eventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(completedEvent, nil).Once()

	svc := NewEventService(eventRepo, schedRepo, partRepo)
	result, err := svc.Complete(context.Background(), entityID, eventID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.EventStatusCompleted, result.Status)

	eventRepo.AssertExpectations(t)
}

func TestEventService_GetByIDWithParticipants(t *testing.T) {
	entityID := uuid.New()
	eventID := uuid.New()

	eventRepo := new(mocks.MockEventRepository)
	schedRepo := new(mocks.MockSchedulerRepository)
	partRepo := new(mocks.MockParticipantRepository)

	event := &domain.Event{ID: eventID, EntityID: entityID, Name: "Event"}
	eventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(event, nil)

	participants := []*domain.Participant{
		{ID: uuid.New(), EventID: eventID, EntityID: entityID, Status: domain.ParticipantStatusConfirmed},
		{ID: uuid.New(), EventID: eventID, EntityID: entityID, Status: domain.ParticipantStatusPending},
	}
	partRepo.On("ListByEvent", mock.Anything, eventID, entityID, 1, 1000).Return(participants, int64(2), nil)

	svc := NewEventService(eventRepo, schedRepo, partRepo)
	result, err := svc.GetByIDWithParticipants(context.Background(), entityID, eventID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Participants, 2)

	eventRepo.AssertExpectations(t)
	partRepo.AssertExpectations(t)
}

// Helper
func eventStrPtr(s string) *string {
	return &s
}
