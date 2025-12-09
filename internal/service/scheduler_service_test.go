package service

import (
	"context"
	"testing"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/testutil/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockNotificationService implements NotificationService for testing
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendConfirmationRequest(ctx context.Context, event *domain.Event, participant *domain.Participant) error {
	args := m.Called(ctx, event, participant)
	return args.Error(0)
}

func (m *MockNotificationService) SendReminder(ctx context.Context, event *domain.Event, participant *domain.Participant) error {
	args := m.Called(ctx, event, participant)
	return args.Error(0)
}

func (m *MockNotificationService) SendLocationRequest(ctx context.Context, event *domain.Event, participant *domain.Participant) error {
	args := m.Called(ctx, event, participant)
	return args.Error(0)
}

func (m *MockNotificationService) SendETAUpdate(ctx context.Context, event *domain.Event, participant *domain.Participant, etaMinutes int) error {
	args := m.Called(ctx, event, participant, etaMinutes)
	return args.Error(0)
}

func (m *MockNotificationService) SendMessage(ctx context.Context, phoneNumber string, message string) error {
	args := m.Called(ctx, phoneNumber, message)
	return args.Error(0)
}

func TestSchedulerService_Create(t *testing.T) {
	logger := zap.NewNop()
	orgID := uuid.New()
	eventID := uuid.New()

	tests := []struct {
		name    string
		setup   func(*mocks.MockSchedulerRepository)
		input   *domain.CreateSchedulerInput
		wantErr bool
	}{
		{
			name: "successful create",
			setup: func(sRepo *mocks.MockSchedulerRepository) {
				sRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			input: &domain.CreateSchedulerInput{
				EventID:     eventID,
				Action:      domain.SchedulerActionReminder,
				ScheduledAt: time.Now().Add(time.Hour),
				MaxRetries:  3,
			},
			wantErr: false,
		},
		{
			name: "create with default max retries",
			setup: func(sRepo *mocks.MockSchedulerRepository) {
				sRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			input: &domain.CreateSchedulerInput{
				EventID:     eventID,
				Action:      domain.SchedulerActionReminder,
				ScheduledAt: time.Now().Add(time.Hour),
				MaxRetries:  0, // Should default to 3
			},
			wantErr: false,
		},
		{
			name: "create fails",
			setup: func(sRepo *mocks.MockSchedulerRepository) {
				sRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			input: &domain.CreateSchedulerInput{
				EventID:     eventID,
				Action:      domain.SchedulerActionReminder,
				ScheduledAt: time.Now().Add(time.Hour),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSchedulerRepo := new(mocks.MockSchedulerRepository)
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			mockNotificationSvc := new(MockNotificationService)

			tt.setup(mockSchedulerRepo)

			svc := NewSchedulerService(mockSchedulerRepo, mockParticipantRepo, mockEventRepo, mockNotificationSvc, logger)
			result, err := svc.Create(context.Background(), tt.input, orgID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, domain.SchedulerStatusPending, result.Status)
			}
		})
	}
}

func TestSchedulerService_GetByID(t *testing.T) {
	logger := zap.NewNop()
	orgID := uuid.New()
	schedulerID := uuid.New()

	tests := []struct {
		name    string
		setup   func(*mocks.MockSchedulerRepository)
		wantErr bool
	}{
		{
			name: "successful get",
			setup: func(sRepo *mocks.MockSchedulerRepository) {
				sRepo.On("GetByID", mock.Anything, schedulerID, orgID).Return(&domain.Scheduler{
					ID:       schedulerID,
					EntityID: orgID,
					Status:   domain.SchedulerStatusPending,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "not found",
			setup: func(sRepo *mocks.MockSchedulerRepository) {
				sRepo.On("GetByID", mock.Anything, schedulerID, orgID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSchedulerRepo := new(mocks.MockSchedulerRepository)
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			mockNotificationSvc := new(MockNotificationService)

			tt.setup(mockSchedulerRepo)

			svc := NewSchedulerService(mockSchedulerRepo, mockParticipantRepo, mockEventRepo, mockNotificationSvc, logger)
			result, err := svc.GetByID(context.Background(), schedulerID, orgID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestSchedulerService_Cancel(t *testing.T) {
	logger := zap.NewNop()
	orgID := uuid.New()
	schedulerID := uuid.New()

	tests := []struct {
		name    string
		setup   func(*mocks.MockSchedulerRepository)
		wantErr bool
	}{
		{
			name: "successful cancel",
			setup: func(sRepo *mocks.MockSchedulerRepository) {
				sRepo.On("GetByID", mock.Anything, schedulerID, orgID).Return(&domain.Scheduler{
					ID:       schedulerID,
					EntityID: orgID,
					Status:   domain.SchedulerStatusPending,
				}, nil)
				sRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "cancel non-pending scheduler",
			setup: func(sRepo *mocks.MockSchedulerRepository) {
				sRepo.On("GetByID", mock.Anything, schedulerID, orgID).Return(&domain.Scheduler{
					ID:       schedulerID,
					EntityID: orgID,
					Status:   domain.SchedulerStatusProcessed,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "scheduler not found",
			setup: func(sRepo *mocks.MockSchedulerRepository) {
				sRepo.On("GetByID", mock.Anything, schedulerID, orgID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSchedulerRepo := new(mocks.MockSchedulerRepository)
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			mockNotificationSvc := new(MockNotificationService)

			tt.setup(mockSchedulerRepo)

			svc := NewSchedulerService(mockSchedulerRepo, mockParticipantRepo, mockEventRepo, mockNotificationSvc, logger)
			err := svc.Cancel(context.Background(), schedulerID, orgID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSchedulerService_ProcessPendingTasks(t *testing.T) {
	logger := zap.NewNop()
	orgID := uuid.New()
	eventID := uuid.New()
	participantID := uuid.New()
	phone := "+5511999999999"

	tests := []struct {
		name          string
		setup         func(*mocks.MockSchedulerRepository, *mocks.MockEventRepository, *mocks.MockParticipantRepository, *MockNotificationService)
		limit         int
		wantProcessed int
		wantErr       bool
	}{
		{
			name: "no pending tasks",
			setup: func(sRepo *mocks.MockSchedulerRepository, eRepo *mocks.MockEventRepository, pRepo *mocks.MockParticipantRepository, nSvc *MockNotificationService) {
				sRepo.On("ListPending", mock.Anything, mock.Anything, 10).Return([]*domain.Scheduler{}, nil)
			},
			limit:         10,
			wantProcessed: 0,
			wantErr:       false,
		},
		{
			name: "list pending error",
			setup: func(sRepo *mocks.MockSchedulerRepository, eRepo *mocks.MockEventRepository, pRepo *mocks.MockParticipantRepository, nSvc *MockNotificationService) {
				sRepo.On("ListPending", mock.Anything, mock.Anything, 10).Return(nil, assert.AnError)
			},
			limit:         10,
			wantProcessed: 0,
			wantErr:       true,
		},
		{
			name: "process confirmation action successfully",
			setup: func(sRepo *mocks.MockSchedulerRepository, eRepo *mocks.MockEventRepository, pRepo *mocks.MockParticipantRepository, nSvc *MockNotificationService) {
				tasks := []*domain.Scheduler{
					{
						ID:       uuid.New(),
						EntityID: orgID,
						EventID:  eventID,
						Action:   domain.SchedulerActionConfirmation,
						Status:   domain.SchedulerStatusPending,
					},
				}
				sRepo.On("ListPending", mock.Anything, mock.Anything, 10).Return(tasks, nil)
				eRepo.On("GetByID", mock.Anything, eventID, orgID).Return(&domain.Event{
					ID:        eventID,
					EntityID:  orgID,
					Name:      "Test Event",
					StartTime: time.Now().Add(24 * time.Hour),
				}, nil)
				pRepo.On("ListByEvent", mock.Anything, eventID, orgID, 1, 1000).Return([]*domain.Participant{
					{
						ID:       participantID,
						EntityID: orgID,
						EventID:  eventID,
						Status:   domain.ParticipantStatusPending,
						Entity: &domain.Entity{
							ID:          uuid.New(),
							Name:        "John Doe",
							PhoneNumber: &phone,
						},
					},
				}, int64(1), nil)
				nSvc.On("SendConfirmationRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				sRepo.On("MarkAsProcessed", mock.Anything, mock.Anything, orgID).Return(nil)
			},
			limit:         10,
			wantProcessed: 1,
			wantErr:       false,
		},
		{
			name: "process reminder action successfully",
			setup: func(sRepo *mocks.MockSchedulerRepository, eRepo *mocks.MockEventRepository, pRepo *mocks.MockParticipantRepository, nSvc *MockNotificationService) {
				tasks := []*domain.Scheduler{
					{
						ID:       uuid.New(),
						EntityID: orgID,
						EventID:  eventID,
						Action:   domain.SchedulerActionReminder,
						Status:   domain.SchedulerStatusPending,
					},
				}
				sRepo.On("ListPending", mock.Anything, mock.Anything, 10).Return(tasks, nil)
				eRepo.On("GetByID", mock.Anything, eventID, orgID).Return(&domain.Event{
					ID:        eventID,
					EntityID:  orgID,
					Name:      "Test Event",
					StartTime: time.Now().Add(1 * time.Hour),
				}, nil)
				pRepo.On("ListByEvent", mock.Anything, eventID, orgID, 1, 1000).Return([]*domain.Participant{
					{
						ID:       participantID,
						EntityID: orgID,
						EventID:  eventID,
						Status:   domain.ParticipantStatusConfirmed,
						Entity: &domain.Entity{
							ID:          uuid.New(),
							Name:        "John Doe",
							PhoneNumber: &phone,
						},
					},
				}, int64(1), nil)
				nSvc.On("SendReminder", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				sRepo.On("MarkAsProcessed", mock.Anything, mock.Anything, orgID).Return(nil)
			},
			limit:         10,
			wantProcessed: 1,
			wantErr:       false,
		},
		{
			name: "process closure action successfully",
			setup: func(sRepo *mocks.MockSchedulerRepository, eRepo *mocks.MockEventRepository, pRepo *mocks.MockParticipantRepository, nSvc *MockNotificationService) {
				tasks := []*domain.Scheduler{
					{
						ID:       uuid.New(),
						EntityID: orgID,
						EventID:  eventID,
						Action:   domain.SchedulerActionClosure,
						Status:   domain.SchedulerStatusPending,
					},
				}
				sRepo.On("ListPending", mock.Anything, mock.Anything, 10).Return(tasks, nil)
				eRepo.On("Update", mock.Anything, eventID, orgID, mock.Anything).Return(nil)
				sRepo.On("MarkAsProcessed", mock.Anything, mock.Anything, orgID).Return(nil)
			},
			limit:         10,
			wantProcessed: 1,
			wantErr:       false,
		},
		{
			name: "process location action successfully",
			setup: func(sRepo *mocks.MockSchedulerRepository, eRepo *mocks.MockEventRepository, pRepo *mocks.MockParticipantRepository, nSvc *MockNotificationService) {
				tasks := []*domain.Scheduler{
					{
						ID:       uuid.New(),
						EntityID: orgID,
						EventID:  eventID,
						Action:   domain.SchedulerActionLocation,
						Status:   domain.SchedulerStatusPending,
					},
				}
				sRepo.On("ListPending", mock.Anything, mock.Anything, 10).Return(tasks, nil)
				eRepo.On("GetByID", mock.Anything, eventID, orgID).Return(&domain.Event{
					ID:        eventID,
					EntityID:  orgID,
					Name:      "Test Event",
					StartTime: time.Now().Add(30 * time.Minute),
				}, nil)
				pRepo.On("ListByEvent", mock.Anything, eventID, orgID, 1, 1000).Return([]*domain.Participant{
					{
						ID:       participantID,
						EntityID: orgID,
						EventID:  eventID,
						Status:   domain.ParticipantStatusConfirmed,
						Entity: &domain.Entity{
							ID:          uuid.New(),
							Name:        "John Doe",
							PhoneNumber: &phone,
						},
					},
				}, int64(1), nil)
				nSvc.On("SendLocationRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				sRepo.On("MarkAsProcessed", mock.Anything, mock.Anything, orgID).Return(nil)
			},
			limit:         10,
			wantProcessed: 1,
			wantErr:       false,
		},
		{
			name: "unknown action is processed without error",
			setup: func(sRepo *mocks.MockSchedulerRepository, eRepo *mocks.MockEventRepository, pRepo *mocks.MockParticipantRepository, nSvc *MockNotificationService) {
				tasks := []*domain.Scheduler{
					{
						ID:       uuid.New(),
						EntityID: orgID,
						EventID:  eventID,
						Action:   domain.SchedulerAction("unknown"),
						Status:   domain.SchedulerStatusPending,
					},
				}
				sRepo.On("ListPending", mock.Anything, mock.Anything, 10).Return(tasks, nil)
				sRepo.On("MarkAsProcessed", mock.Anything, mock.Anything, orgID).Return(nil)
			},
			limit:         10,
			wantProcessed: 1,
			wantErr:       false,
		},
		{
			name: "task processing fails and increments retries",
			setup: func(sRepo *mocks.MockSchedulerRepository, eRepo *mocks.MockEventRepository, pRepo *mocks.MockParticipantRepository, nSvc *MockNotificationService) {
				taskID := uuid.New()
				tasks := []*domain.Scheduler{
					{
						ID:         taskID,
						EntityID:   orgID,
						EventID:    eventID,
						Action:     domain.SchedulerActionConfirmation,
						Status:     domain.SchedulerStatusPending,
						Retries:    0,
						MaxRetries: 3,
					},
				}
				sRepo.On("ListPending", mock.Anything, mock.Anything, 10).Return(tasks, nil)
				eRepo.On("GetByID", mock.Anything, eventID, orgID).Return(nil, assert.AnError)
				sRepo.On("IncrementRetries", mock.Anything, taskID, orgID).Return(nil)
			},
			limit:         10,
			wantProcessed: 0,
			wantErr:       false,
		},
		{
			name: "task processing fails and marks as failed when max retries exceeded",
			setup: func(sRepo *mocks.MockSchedulerRepository, eRepo *mocks.MockEventRepository, pRepo *mocks.MockParticipantRepository, nSvc *MockNotificationService) {
				taskID := uuid.New()
				tasks := []*domain.Scheduler{
					{
						ID:         taskID,
						EntityID:   orgID,
						EventID:    eventID,
						Action:     domain.SchedulerActionConfirmation,
						Status:     domain.SchedulerStatusPending,
						Retries:    2,
						MaxRetries: 3,
					},
				}
				sRepo.On("ListPending", mock.Anything, mock.Anything, 10).Return(tasks, nil)
				eRepo.On("GetByID", mock.Anything, eventID, orgID).Return(nil, assert.AnError)
				sRepo.On("IncrementRetries", mock.Anything, taskID, orgID).Return(nil)
				sRepo.On("MarkAsFailed", mock.Anything, taskID, orgID, mock.Anything).Return(nil)
			},
			limit:         10,
			wantProcessed: 0,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSchedulerRepo := new(mocks.MockSchedulerRepository)
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			mockNotificationSvc := new(MockNotificationService)

			tt.setup(mockSchedulerRepo, mockEventRepo, mockParticipantRepo, mockNotificationSvc)

			svc := NewSchedulerService(mockSchedulerRepo, mockParticipantRepo, mockEventRepo, mockNotificationSvc, logger)
			processed, err := svc.ProcessPendingTasks(context.Background(), tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantProcessed, processed)
			}
		})
	}
}

func TestSchedulerService_ProcessConfirmation_SkipsNonPending(t *testing.T) {
	logger := zap.NewNop()
	orgID := uuid.New()
	eventID := uuid.New()
	phone := "+5511999999999"

	mockSchedulerRepo := new(mocks.MockSchedulerRepository)
	mockParticipantRepo := new(mocks.MockParticipantRepository)
	mockEventRepo := new(mocks.MockEventRepository)
	mockNotificationSvc := new(MockNotificationService)

	tasks := []*domain.Scheduler{
		{
			ID:       uuid.New(),
			EntityID: orgID,
			EventID:  eventID,
			Action:   domain.SchedulerActionConfirmation,
			Status:   domain.SchedulerStatusPending,
		},
	}
	mockSchedulerRepo.On("ListPending", mock.Anything, mock.Anything, 10).Return(tasks, nil)
	mockEventRepo.On("GetByID", mock.Anything, eventID, orgID).Return(&domain.Event{
		ID:        eventID,
		EntityID:  orgID,
		Name:      "Test Event",
		StartTime: time.Now().Add(24 * time.Hour),
	}, nil)
	// Return participant with Confirmed status (should be skipped)
	mockParticipantRepo.On("ListByEvent", mock.Anything, eventID, orgID, 1, 1000).Return([]*domain.Participant{
		{
			ID:       uuid.New(),
			EntityID: orgID,
			EventID:  eventID,
			Status:   domain.ParticipantStatusConfirmed, // Already confirmed
			Entity: &domain.Entity{
				ID:          uuid.New(),
				Name:        "John Doe",
				PhoneNumber: &phone,
			},
		},
	}, int64(1), nil)
	mockSchedulerRepo.On("MarkAsProcessed", mock.Anything, mock.Anything, orgID).Return(nil)

	svc := NewSchedulerService(mockSchedulerRepo, mockParticipantRepo, mockEventRepo, mockNotificationSvc, logger)
	processed, err := svc.ProcessPendingTasks(context.Background(), 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, processed)
	// SendConfirmationRequest should NOT be called
	mockNotificationSvc.AssertNotCalled(t, "SendConfirmationRequest")
}

func TestSchedulerService_ProcessReminder_SkipsNonConfirmed(t *testing.T) {
	logger := zap.NewNop()
	orgID := uuid.New()
	eventID := uuid.New()
	phone := "+5511999999999"

	mockSchedulerRepo := new(mocks.MockSchedulerRepository)
	mockParticipantRepo := new(mocks.MockParticipantRepository)
	mockEventRepo := new(mocks.MockEventRepository)
	mockNotificationSvc := new(MockNotificationService)

	tasks := []*domain.Scheduler{
		{
			ID:       uuid.New(),
			EntityID: orgID,
			EventID:  eventID,
			Action:   domain.SchedulerActionReminder,
			Status:   domain.SchedulerStatusPending,
		},
	}
	mockSchedulerRepo.On("ListPending", mock.Anything, mock.Anything, 10).Return(tasks, nil)
	mockEventRepo.On("GetByID", mock.Anything, eventID, orgID).Return(&domain.Event{
		ID:        eventID,
		EntityID:  orgID,
		Name:      "Test Event",
		StartTime: time.Now().Add(1 * time.Hour),
	}, nil)
	// Return participant with Pending status (should be skipped for reminder)
	mockParticipantRepo.On("ListByEvent", mock.Anything, eventID, orgID, 1, 1000).Return([]*domain.Participant{
		{
			ID:       uuid.New(),
			EntityID: orgID,
			EventID:  eventID,
			Status:   domain.ParticipantStatusPending, // Not confirmed
			Entity: &domain.Entity{
				ID:          uuid.New(),
				Name:        "John Doe",
				PhoneNumber: &phone,
			},
		},
	}, int64(1), nil)
	mockSchedulerRepo.On("MarkAsProcessed", mock.Anything, mock.Anything, orgID).Return(nil)

	svc := NewSchedulerService(mockSchedulerRepo, mockParticipantRepo, mockEventRepo, mockNotificationSvc, logger)
	processed, err := svc.ProcessPendingTasks(context.Background(), 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, processed)
	// SendReminder should NOT be called
	mockNotificationSvc.AssertNotCalled(t, "SendReminder")
}

func TestSchedulerService_ProcessConfirmation_NotificationError(t *testing.T) {
	logger := zap.NewNop()
	orgID := uuid.New()
	eventID := uuid.New()
	phone := "+5511999999999"

	mockSchedulerRepo := new(mocks.MockSchedulerRepository)
	mockParticipantRepo := new(mocks.MockParticipantRepository)
	mockEventRepo := new(mocks.MockEventRepository)
	mockNotificationSvc := new(MockNotificationService)

	tasks := []*domain.Scheduler{
		{
			ID:       uuid.New(),
			EntityID: orgID,
			EventID:  eventID,
			Action:   domain.SchedulerActionConfirmation,
			Status:   domain.SchedulerStatusPending,
		},
	}
	mockSchedulerRepo.On("ListPending", mock.Anything, mock.Anything, 10).Return(tasks, nil)
	mockEventRepo.On("GetByID", mock.Anything, eventID, orgID).Return(&domain.Event{
		ID:        eventID,
		EntityID:  orgID,
		Name:      "Test Event",
		StartTime: time.Now().Add(24 * time.Hour),
	}, nil)
	mockParticipantRepo.On("ListByEvent", mock.Anything, eventID, orgID, 1, 1000).Return([]*domain.Participant{
		{
			ID:       uuid.New(),
			EntityID: orgID,
			EventID:  eventID,
			Status:   domain.ParticipantStatusPending,
			Entity: &domain.Entity{
				ID:          uuid.New(),
				Name:        "John Doe",
				PhoneNumber: &phone,
			},
		},
	}, int64(1), nil)
	// Notification fails but processing continues
	mockNotificationSvc.On("SendConfirmationRequest", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)
	mockSchedulerRepo.On("MarkAsProcessed", mock.Anything, mock.Anything, orgID).Return(nil)

	svc := NewSchedulerService(mockSchedulerRepo, mockParticipantRepo, mockEventRepo, mockNotificationSvc, logger)
	processed, err := svc.ProcessPendingTasks(context.Background(), 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, processed)
}
