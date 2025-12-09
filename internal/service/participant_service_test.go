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

func TestParticipantService_Create(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*mocks.MockParticipantRepository, *mocks.MockEventRepository)
		req         *dto.CreateParticipantRequest
		wantErr     bool
		errContains string
	}{
		{
			name: "successful create",
			setup: func(pRepo *mocks.MockParticipantRepository, eRepo *mocks.MockEventRepository) {
				eventID := uuid.New()
				entID := uuid.New()
				eRepo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).Return(&domain.Event{
					ID:       eventID,
					EntityID: entID,
				}, nil)
				pRepo.On("GetByPhoneNumber", mock.Anything, "5511999999999", mock.Anything, mock.Anything).Return(nil, domain.ErrNotFound)
				pRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			req: &dto.CreateParticipantRequest{
				PhoneNumber: "5511999999999",
			},
			wantErr: false,
		},
		{
			name: "event not found",
			setup: func(pRepo *mocks.MockParticipantRepository, eRepo *mocks.MockEventRepository) {
				eRepo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).Return(nil, domain.ErrNotFound)
			},
			req: &dto.CreateParticipantRequest{
				PhoneNumber: "5511999999999",
			},
			wantErr:     true,
			errContains: "event not found",
		},
		{
			name: "duplicate phone number",
			setup: func(pRepo *mocks.MockParticipantRepository, eRepo *mocks.MockEventRepository) {
				eRepo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).Return(&domain.Event{
					ID: uuid.New(),
				}, nil)
				pRepo.On("GetByPhoneNumber", mock.Anything, "5511999999999", mock.Anything, mock.Anything).Return(&domain.Participant{}, nil)
			},
			req: &dto.CreateParticipantRequest{
				PhoneNumber: "5511999999999",
			},
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)

			tt.setup(mockParticipantRepo, mockEventRepo)

			svc := NewParticipantService(mockParticipantRepo, mockEventRepo)
			result, err := svc.Create(context.Background(), uuid.New(), uuid.New(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestParticipantService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mocks.MockParticipantRepository)
		wantErr bool
	}{
		{
			name: "successful get",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).Return(&domain.Participant{
					ID:       uuid.New(),
					EntityID: uuid.New(),
					EventID:  uuid.New(),
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "not found",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			tt.setup(mockParticipantRepo)

			svc := NewParticipantService(mockParticipantRepo, mockEventRepo)
			result, err := svc.GetByID(context.Background(), uuid.New(), uuid.New())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestParticipantService_Update(t *testing.T) {
	participantID := uuid.New()
	entID := uuid.New()

	tests := []struct {
		name    string
		setup   func(*mocks.MockParticipantRepository)
		req     *dto.UpdateParticipantRequest
		wantErr bool
	}{
		{
			name: "successful update",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetByID", mock.Anything, participantID, entID).Return(&domain.Participant{
					ID:       participantID,
					EntityID: entID,
					EventID:  uuid.New(),
				}, nil)
				pRepo.On("Update", mock.Anything, participantID, entID, mock.Anything).Return(nil)
			},
			req: &dto.UpdateParticipantRequest{
				Name: participantStrPtr("Updated Name"),
			},
			wantErr: false,
		},
		{
			name: "update with status change to confirmed",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetByID", mock.Anything, participantID, entID).Return(&domain.Participant{
					ID:       participantID,
					EntityID: entID,
					EventID:  uuid.New(),
				}, nil)
				pRepo.On("Update", mock.Anything, participantID, entID, mock.Anything).Return(nil)
			},
			req: &dto.UpdateParticipantRequest{
				Status: participantStatusPtr(domain.ParticipantStatusConfirmed),
			},
			wantErr: false,
		},
		{
			name: "participant not found",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetByID", mock.Anything, participantID, entID).Return(nil, domain.ErrNotFound)
			},
			req: &dto.UpdateParticipantRequest{
				Name: participantStrPtr("Updated Name"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			tt.setup(mockParticipantRepo)

			svc := NewParticipantService(mockParticipantRepo, mockEventRepo)
			result, err := svc.Update(context.Background(), entID, participantID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestParticipantService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mocks.MockParticipantRepository)
		wantErr bool
	}{
		{
			name: "successful delete",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "not found",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			tt.setup(mockParticipantRepo)

			svc := NewParticipantService(mockParticipantRepo, mockEventRepo)
			err := svc.Delete(context.Background(), uuid.New(), uuid.New())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParticipantService_ListByEvent(t *testing.T) {
	eventID := uuid.New()
	entID := uuid.New()

	tests := []struct {
		name    string
		setup   func(*mocks.MockParticipantRepository, *mocks.MockEventRepository)
		wantErr bool
	}{
		{
			name: "successful list",
			setup: func(pRepo *mocks.MockParticipantRepository, eRepo *mocks.MockEventRepository) {
				eRepo.On("GetByID", mock.Anything, eventID, entID).Return(&domain.Event{
					ID: eventID,
				}, nil)
				pRepo.On("ListByEvent", mock.Anything, eventID, entID, 1, 10).Return([]*domain.Participant{
					{ID: uuid.New(), EntityID: entID, EventID: eventID},
				}, int64(1), nil)
			},
			wantErr: false,
		},
		{
			name: "event not found",
			setup: func(pRepo *mocks.MockParticipantRepository, eRepo *mocks.MockEventRepository) {
				eRepo.On("GetByID", mock.Anything, eventID, entID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			tt.setup(mockParticipantRepo, mockEventRepo)

			svc := NewParticipantService(mockParticipantRepo, mockEventRepo)
			result, total, err := svc.ListByEvent(context.Background(), entID, eventID, 1, 10)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, int64(1), total)
			}
		})
	}
}

func TestParticipantService_UpdateStatus(t *testing.T) {
	participantID := uuid.New()
	entID := uuid.New()

	tests := []struct {
		name    string
		setup   func(*mocks.MockParticipantRepository)
		status  domain.ParticipantStatus
		wantErr bool
	}{
		{
			name: "successful status update",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("UpdateStatus", mock.Anything, participantID, entID, domain.ParticipantStatusConfirmed).Return(nil)
			},
			status:  domain.ParticipantStatusConfirmed,
			wantErr: false,
		},
		{
			name: "status update error",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("UpdateStatus", mock.Anything, participantID, entID, domain.ParticipantStatusDenied).Return(domain.ErrNotFound)
			},
			status:  domain.ParticipantStatusDenied,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			tt.setup(mockParticipantRepo)

			svc := NewParticipantService(mockParticipantRepo, mockEventRepo)
			err := svc.UpdateStatus(context.Background(), entID, participantID, tt.status)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParticipantService_ConfirmParticipant(t *testing.T) {
	participantID := uuid.New()
	entID := uuid.New()
	eventID := uuid.New()

	tests := []struct {
		name    string
		setup   func(*mocks.MockParticipantRepository)
		wantErr bool
	}{
		{
			name: "successful confirm",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetByID", mock.Anything, participantID, entID).Return(&domain.Participant{
					ID:       participantID,
					EntityID: entID,
					EventID:  eventID,
					Status:   domain.ParticipantStatusPending,
				}, nil)
				pRepo.On("Update", mock.Anything, participantID, entID, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "participant not found",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetByID", mock.Anything, participantID, entID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			tt.setup(mockParticipantRepo)

			svc := NewParticipantService(mockParticipantRepo, mockEventRepo)
			result, err := svc.ConfirmParticipant(context.Background(), entID, participantID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestParticipantService_CheckInParticipant(t *testing.T) {
	participantID := uuid.New()
	entID := uuid.New()
	eventID := uuid.New()

	tests := []struct {
		name    string
		setup   func(*mocks.MockParticipantRepository)
		wantErr bool
	}{
		{
			name: "successful check-in",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetByID", mock.Anything, participantID, entID).Return(&domain.Participant{
					ID:       participantID,
					EntityID: entID,
					EventID:  eventID,
					Status:   domain.ParticipantStatusConfirmed,
				}, nil)
				pRepo.On("Update", mock.Anything, participantID, entID, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "participant not found",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetByID", mock.Anything, participantID, entID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			tt.setup(mockParticipantRepo)

			svc := NewParticipantService(mockParticipantRepo, mockEventRepo)
			result, err := svc.CheckInParticipant(context.Background(), entID, participantID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestParticipantService_BatchCreate(t *testing.T) {
	entID := uuid.New()
	eventID := uuid.New()

	tests := []struct {
		name          string
		setup         func(*mocks.MockParticipantRepository, *mocks.MockEventRepository)
		req           *dto.BatchCreateParticipantsRequest
		wantResponses int
		wantErrors    int
	}{
		{
			name: "successful batch create",
			setup: func(pRepo *mocks.MockParticipantRepository, eRepo *mocks.MockEventRepository) {
				eRepo.On("GetByID", mock.Anything, eventID, entID).Return(&domain.Event{
					ID:       eventID,
					EntityID: entID,
				}, nil)
				pRepo.On("GetByPhoneNumber", mock.Anything, "5511999999991", eventID, entID).Return(nil, domain.ErrNotFound)
				pRepo.On("GetByPhoneNumber", mock.Anything, "5511999999992", eventID, entID).Return(nil, domain.ErrNotFound)
				pRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			req: &dto.BatchCreateParticipantsRequest{
				Participants: []dto.CreateParticipantRequest{
					{PhoneNumber: "5511999999991"},
					{PhoneNumber: "5511999999992"},
				},
			},
			wantResponses: 2,
			wantErrors:    0,
		},
		{
			name: "event not found",
			setup: func(pRepo *mocks.MockParticipantRepository, eRepo *mocks.MockEventRepository) {
				eRepo.On("GetByID", mock.Anything, eventID, entID).Return(nil, domain.ErrNotFound)
			},
			req: &dto.BatchCreateParticipantsRequest{
				Participants: []dto.CreateParticipantRequest{
					{PhoneNumber: "5511999999991"},
				},
			},
			wantResponses: 0,
			wantErrors:    1,
		},
		{
			name: "partial failure - one duplicate",
			setup: func(pRepo *mocks.MockParticipantRepository, eRepo *mocks.MockEventRepository) {
				eRepo.On("GetByID", mock.Anything, eventID, entID).Return(&domain.Event{
					ID:       eventID,
					EntityID: entID,
				}, nil)
				pRepo.On("GetByPhoneNumber", mock.Anything, "5511999999991", eventID, entID).Return(nil, domain.ErrNotFound)
				pRepo.On("GetByPhoneNumber", mock.Anything, "5511999999992", eventID, entID).Return(&domain.Participant{}, nil)
				pRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			req: &dto.BatchCreateParticipantsRequest{
				Participants: []dto.CreateParticipantRequest{
					{PhoneNumber: "5511999999991"},
					{PhoneNumber: "5511999999992"},
				},
			},
			wantResponses: 1,
			wantErrors:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			tt.setup(mockParticipantRepo, mockEventRepo)

			svc := NewParticipantService(mockParticipantRepo, mockEventRepo)
			responses, errors := svc.BatchCreate(context.Background(), entID, eventID, tt.req)

			assert.Len(t, responses, tt.wantResponses)
			assert.Len(t, errors, tt.wantErrors)
		})
	}
}

func TestParticipantService_GetByPhoneNumber(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*mocks.MockParticipantRepository)
		phoneNumber string
		wantErr     bool
	}{
		{
			name: "successful get by phone",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetActiveByPhoneNumber", mock.Anything, "5511999999999").Return(&domain.Participant{
					ID:       uuid.New(),
					EntityID: uuid.New(),
					EventID:  uuid.New(),
				}, nil)
			},
			phoneNumber: "5511999999999",
			wantErr:     false,
		},
		{
			name: "not found",
			setup: func(pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetActiveByPhoneNumber", mock.Anything, "5511888888888").Return(nil, domain.ErrNotFound)
			},
			phoneNumber: "5511888888888",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)
			tt.setup(mockParticipantRepo)

			svc := NewParticipantService(mockParticipantRepo, mockEventRepo)
			result, err := svc.GetByPhoneNumber(context.Background(), tt.phoneNumber)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// Helper functions
func participantStrPtr(s string) *string {
	return &s
}

func participantStatusPtr(s domain.ParticipantStatus) *domain.ParticipantStatus {
	return &s
}

func participantTimePtr(t time.Time) *time.Time {
	return &t
}
