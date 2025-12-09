package service

import (
	"context"
	"testing"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewNotificationService(t *testing.T) {
	logger := zap.NewNop()
	svc := NewNotificationService(nil, logger)
	assert.NotNil(t, svc)
}

func TestNotificationService_SendConfirmationRequest(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	phone := "+5511999999999"

	tests := []struct {
		name        string
		event       *domain.Event
		participant *domain.Participant
		wantErr     bool
	}{
		{
			name: "successful send with phone",
			event: &domain.Event{
				ID:        uuid.New(),
				Name:      "Test Event",
				StartTime: time.Now().Add(24 * time.Hour),
			},
			participant: &domain.Participant{
				ID: uuid.New(),
				Entity: &domain.Entity{
					ID:          uuid.New(),
					Name:        "John Doe",
					PhoneNumber: &phone,
				},
			},
			wantErr: false,
		},
		{
			name: "participant without entity",
			event: &domain.Event{
				ID:   uuid.New(),
				Name: "Test Event",
			},
			participant: &domain.Participant{
				ID:     uuid.New(),
				Entity: nil,
			},
			wantErr: false, // No error, just skips
		},
		{
			name: "participant entity without phone",
			event: &domain.Event{
				ID:   uuid.New(),
				Name: "Test Event",
			},
			participant: &domain.Participant{
				ID: uuid.New(),
				Entity: &domain.Entity{
					ID:          uuid.New(),
					Name:        "John Doe",
					PhoneNumber: nil,
				},
			},
			wantErr: false, // No error, just skips
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create service with nil WhatsApp client (skips actual sending)
			svc := NewNotificationService(nil, logger)
			err := svc.SendConfirmationRequest(ctx, tt.event, tt.participant)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNotificationService_SendReminder(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	phone := "+5511999999999"
	address := "Av. Paulista, 1000"

	tests := []struct {
		name        string
		event       *domain.Event
		participant *domain.Participant
		wantErr     bool
	}{
		{
			name: "successful send with address",
			event: &domain.Event{
				ID:              uuid.New(),
				Name:            "Test Event",
				StartTime:       time.Now().Add(1 * time.Hour),
				LocationAddress: &address,
			},
			participant: &domain.Participant{
				ID: uuid.New(),
				Entity: &domain.Entity{
					ID:          uuid.New(),
					Name:        "John Doe",
					PhoneNumber: &phone,
				},
			},
			wantErr: false,
		},
		{
			name: "successful send with coordinates only",
			event: &domain.Event{
				ID:          uuid.New(),
				Name:        "Test Event",
				StartTime:   time.Now().Add(1 * time.Hour),
				LocationLat: -23.5505,
				LocationLng: -46.6333,
			},
			participant: &domain.Participant{
				ID: uuid.New(),
				Entity: &domain.Entity{
					ID:          uuid.New(),
					Name:        "John Doe",
					PhoneNumber: &phone,
				},
			},
			wantErr: false,
		},
		{
			name: "participant without phone",
			event: &domain.Event{
				ID:   uuid.New(),
				Name: "Test Event",
			},
			participant: &domain.Participant{
				ID: uuid.New(),
				Entity: &domain.Entity{
					ID:   uuid.New(),
					Name: "John Doe",
				},
			},
			wantErr: false, // No error, just skips
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewNotificationService(nil, logger)
			err := svc.SendReminder(ctx, tt.event, tt.participant)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNotificationService_SendLocationRequest(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	phone := "+5511999999999"

	tests := []struct {
		name        string
		event       *domain.Event
		participant *domain.Participant
		wantErr     bool
	}{
		{
			name: "successful send",
			event: &domain.Event{
				ID:   uuid.New(),
				Name: "Test Event",
			},
			participant: &domain.Participant{
				ID: uuid.New(),
				Entity: &domain.Entity{
					ID:          uuid.New(),
					Name:        "John Doe",
					PhoneNumber: &phone,
				},
			},
			wantErr: false,
		},
		{
			name: "participant without entity",
			event: &domain.Event{
				ID:   uuid.New(),
				Name: "Test Event",
			},
			participant: &domain.Participant{
				ID:     uuid.New(),
				Entity: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewNotificationService(nil, logger)
			err := svc.SendLocationRequest(ctx, tt.event, tt.participant)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNotificationService_SendETAUpdate(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	tests := []struct {
		name       string
		etaMinutes int
	}{
		{name: "less than 5 minutes", etaMinutes: 3},
		{name: "5 to 60 minutes", etaMinutes: 30},
		{name: "more than 60 minutes", etaMinutes: 90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewNotificationService(nil, logger)
			event := &domain.Event{ID: uuid.New(), Name: "Test"}
			participant := &domain.Participant{ID: uuid.New()}

			err := svc.SendETAUpdate(ctx, event, participant, tt.etaMinutes)
			assert.NoError(t, err)
		})
	}
}

func TestNotificationService_SendMessage(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	t.Run("nil whatsapp client skips sending", func(t *testing.T) {
		svc := NewNotificationService(nil, logger)
		err := svc.SendMessage(ctx, "+5511999999999", "Test message")
		assert.NoError(t, err)
	})
}

func TestGetLocationAddress(t *testing.T) {
	tests := []struct {
		name     string
		event    *domain.Event
		expected string
	}{
		{
			name: "with address",
			event: &domain.Event{
				LocationAddress: notifStrPtr("Av. Paulista, 1000"),
				LocationLat:     -23.5505,
				LocationLng:     -46.6333,
			},
			expected: "Av. Paulista, 1000",
		},
		{
			name: "with empty address",
			event: &domain.Event{
				LocationAddress: notifStrPtr(""),
				LocationLat:     -23.5505,
				LocationLng:     -46.6333,
			},
			expected: "-23.550500, -46.633300",
		},
		{
			name: "with nil address",
			event: &domain.Event{
				LocationLat: -23.5505,
				LocationLng: -46.6333,
			},
			expected: "-23.550500, -46.633300",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLocationAddress(tt.event)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func notifStrPtr(s string) *string {
	return &s
}
