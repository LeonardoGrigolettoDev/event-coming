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
	"go.uber.org/zap"
)

func TestLocationService_CreateLocation(t *testing.T) {
	logger := zap.NewNop()
	participantID := uuid.New()
	entityID := uuid.New()
	eventID := uuid.New()

	tests := []struct {
		name    string
		setup   func(*mocks.MockLocationRepository, *mocks.MockParticipantRepository, *mocks.MockEventRepository)
		req     *dto.CreateLocationRequest
		wantErr bool
	}{
		{
			name: "successful create",
			setup: func(lRepo *mocks.MockLocationRepository, pRepo *mocks.MockParticipantRepository, eRepo *mocks.MockEventRepository) {
				pRepo.On("GetByID", mock.Anything, participantID, entityID).Return(&domain.Participant{
					ID:       participantID,
					EventID:  eventID,
					EntityID: entityID,
				}, nil)
				eRepo.On("GetByID", mock.Anything, eventID, entityID).Return(&domain.Event{
					ID: eventID,
				}, nil)
				lRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			req: &dto.CreateLocationRequest{
				Latitude:  -23.5505,
				Longitude: -46.6333,
			},
			wantErr: false,
		},
		{
			name: "participant not found",
			setup: func(lRepo *mocks.MockLocationRepository, pRepo *mocks.MockParticipantRepository, eRepo *mocks.MockEventRepository) {
				pRepo.On("GetByID", mock.Anything, participantID, entityID).Return(nil, domain.ErrNotFound)
			},
			req: &dto.CreateLocationRequest{
				Latitude:  -23.5505,
				Longitude: -46.6333,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLocationRepo := new(mocks.MockLocationRepository)
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)

			tt.setup(mockLocationRepo, mockParticipantRepo, mockEventRepo)

			svc := NewLocationService(mockLocationRepo, mockParticipantRepo, mockEventRepo, nil, logger)
			result, err := svc.CreateLocation(context.Background(), participantID, entityID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestLocationService_GetLatestLocation(t *testing.T) {
	logger := zap.NewNop()
	participantID := uuid.New()
	entityID := uuid.New()
	eventID := uuid.New()

	tests := []struct {
		name    string
		setup   func(*mocks.MockLocationRepository, *mocks.MockParticipantRepository)
		wantErr bool
	}{
		{
			name: "successful get from database",
			setup: func(lRepo *mocks.MockLocationRepository, pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetByID", mock.Anything, participantID, entityID).Return(&domain.Participant{
					ID:       participantID,
					EventID:  eventID,
					EntityID: entityID,
				}, nil)
				lRepo.On("GetLatestByParticipant", mock.Anything, participantID, entityID).Return(&domain.Location{
					ID:            uuid.New(),
					ParticipantID: participantID,
					Latitude:      -23.5505,
					Longitude:     -46.6333,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "participant not found",
			setup: func(lRepo *mocks.MockLocationRepository, pRepo *mocks.MockParticipantRepository) {
				pRepo.On("GetByID", mock.Anything, participantID, entityID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLocationRepo := new(mocks.MockLocationRepository)
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)

			tt.setup(mockLocationRepo, mockParticipantRepo)

			svc := NewLocationService(mockLocationRepo, mockParticipantRepo, mockEventRepo, nil, logger)
			result, err := svc.GetLatestLocation(context.Background(), participantID, entityID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestLocationService_GetLocationHistory(t *testing.T) {
	logger := zap.NewNop()
	participantID := uuid.New()
	entityID := uuid.New()
	from := time.Now().Add(-time.Hour)
	to := time.Now()

	tests := []struct {
		name    string
		setup   func(*mocks.MockLocationRepository)
		wantErr bool
	}{
		{
			name: "successful get history",
			setup: func(lRepo *mocks.MockLocationRepository) {
				lRepo.On("GetHistory", mock.Anything, participantID, entityID, from, to).Return([]*domain.Location{
					{ID: uuid.New(), ParticipantID: participantID, Latitude: -23.5505, Longitude: -46.6333},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "empty history",
			setup: func(lRepo *mocks.MockLocationRepository) {
				lRepo.On("GetHistory", mock.Anything, participantID, entityID, from, to).Return([]*domain.Location{}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLocationRepo := new(mocks.MockLocationRepository)
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)

			tt.setup(mockLocationRepo)

			svc := NewLocationService(mockLocationRepo, mockParticipantRepo, mockEventRepo, nil, logger)
			result, err := svc.GetLocationHistory(context.Background(), participantID, entityID, from, to)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestNewLocationService(t *testing.T) {
	logger := zap.NewNop()
	mockLocationRepo := new(mocks.MockLocationRepository)
	mockParticipantRepo := new(mocks.MockParticipantRepository)
	mockEventRepo := new(mocks.MockEventRepository)

	svc := NewLocationService(mockLocationRepo, mockParticipantRepo, mockEventRepo, nil, logger)
	assert.NotNil(t, svc)
}

func TestLocationService_GetEventLocations(t *testing.T) {
	logger := zap.NewNop()
	entityID := uuid.New()
	eventID := uuid.New()
	participantID := uuid.New()

	tests := []struct {
		name    string
		setup   func(*mocks.MockLocationRepository, *mocks.MockParticipantRepository)
		wantErr bool
	}{
		{
			name: "successful get event locations",
			setup: func(lRepo *mocks.MockLocationRepository, pRepo *mocks.MockParticipantRepository) {
				lRepo.On("GetLatestByEvent", mock.Anything, eventID, entityID).Return([]*domain.Location{
					{
						ID:            uuid.New(),
						ParticipantID: participantID,
						EventID:       eventID,
						Latitude:      -23.5505,
						Longitude:     -46.6333,
						Timestamp:     time.Now(),
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			setup: func(lRepo *mocks.MockLocationRepository, pRepo *mocks.MockParticipantRepository) {
				lRepo.On("GetLatestByEvent", mock.Anything, eventID, entityID).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "empty locations",
			setup: func(lRepo *mocks.MockLocationRepository, pRepo *mocks.MockParticipantRepository) {
				lRepo.On("GetLatestByEvent", mock.Anything, eventID, entityID).Return([]*domain.Location{}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLocationRepo := new(mocks.MockLocationRepository)
			mockParticipantRepo := new(mocks.MockParticipantRepository)
			mockEventRepo := new(mocks.MockEventRepository)

			tt.setup(mockLocationRepo, mockParticipantRepo)

			svc := NewLocationService(mockLocationRepo, mockParticipantRepo, mockEventRepo, nil, logger)
			result, err := svc.GetEventLocations(context.Background(), eventID, entityID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestLocationService_CreateLocation_WithOptionalFields(t *testing.T) {
	logger := zap.NewNop()
	participantID := uuid.New()
	entityID := uuid.New()
	eventID := uuid.New()

	accuracy := 10.0
	speed := 5.0
	heading := 180.0

	mockLocationRepo := new(mocks.MockLocationRepository)
	mockParticipantRepo := new(mocks.MockParticipantRepository)
	mockEventRepo := new(mocks.MockEventRepository)

	mockParticipantRepo.On("GetByID", mock.Anything, participantID, entityID).Return(&domain.Participant{
		ID:       participantID,
		EventID:  eventID,
		EntityID: entityID,
	}, nil)
	mockEventRepo.On("GetByID", mock.Anything, eventID, entityID).Return(&domain.Event{
		ID:       eventID,
		EntityID: entityID,
	}, nil)
	mockLocationRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	svc := NewLocationService(mockLocationRepo, mockParticipantRepo, mockEventRepo, nil, logger)
	result, err := svc.CreateLocation(context.Background(), participantID, entityID, &dto.CreateLocationRequest{
		Latitude:  -23.5505,
		Longitude: -46.6333,
		Accuracy:  &accuracy,
		Speed:     &speed,
		Heading:   &heading,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, -23.5505, result.Latitude)
	assert.Equal(t, -46.6333, result.Longitude)
}

func TestLocationService_GetLatestLocation_NotFound(t *testing.T) {
	logger := zap.NewNop()
	participantID := uuid.New()
	entityID := uuid.New()
	eventID := uuid.New()

	mockLocationRepo := new(mocks.MockLocationRepository)
	mockParticipantRepo := new(mocks.MockParticipantRepository)
	mockEventRepo := new(mocks.MockEventRepository)

	mockParticipantRepo.On("GetByID", mock.Anything, participantID, entityID).Return(&domain.Participant{
		ID:       participantID,
		EventID:  eventID,
		EntityID: entityID,
	}, nil)
	mockLocationRepo.On("GetLatestByParticipant", mock.Anything, participantID, entityID).Return(nil, domain.ErrNotFound)

	svc := NewLocationService(mockLocationRepo, mockParticipantRepo, mockEventRepo, nil, logger)
	result, err := svc.GetLatestLocation(context.Background(), participantID, entityID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestLocationService_GetLocationHistory_Error(t *testing.T) {
	logger := zap.NewNop()
	participantID := uuid.New()
	entityID := uuid.New()
	from := time.Now().Add(-time.Hour)
	to := time.Now()

	mockLocationRepo := new(mocks.MockLocationRepository)
	mockParticipantRepo := new(mocks.MockParticipantRepository)
	mockEventRepo := new(mocks.MockEventRepository)

	mockLocationRepo.On("GetHistory", mock.Anything, participantID, entityID, from, to).Return(nil, assert.AnError)

	svc := NewLocationService(mockLocationRepo, mockParticipantRepo, mockEventRepo, nil, logger)
	result, err := svc.GetLocationHistory(context.Background(), participantID, entityID, from, to)

	assert.Error(t, err)
	assert.Nil(t, result)
}
