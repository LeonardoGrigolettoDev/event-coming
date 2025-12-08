package domain_test

import (
	"testing"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestScheduler_TableName(t *testing.T) {
	scheduler := domain.Scheduler{}
	assert.Equal(t, "schedulers", scheduler.TableName())
}

func TestSchedulerAction_Constants(t *testing.T) {
	assert.Equal(t, domain.SchedulerAction("confirmation"), domain.SchedulerActionConfirmation)
	assert.Equal(t, domain.SchedulerAction("reminder"), domain.SchedulerActionReminder)
	assert.Equal(t, domain.SchedulerAction("closure"), domain.SchedulerActionClosure)
	assert.Equal(t, domain.SchedulerAction("location"), domain.SchedulerActionLocation)
}

func TestSchedulerStatus_Constants(t *testing.T) {
	assert.Equal(t, domain.SchedulerStatus("pending"), domain.SchedulerStatusPending)
	assert.Equal(t, domain.SchedulerStatus("processed"), domain.SchedulerStatusProcessed)
	assert.Equal(t, domain.SchedulerStatus("failed"), domain.SchedulerStatusFailed)
	assert.Equal(t, domain.SchedulerStatus("skipped"), domain.SchedulerStatusSkipped)
}

func TestScheduler_Fields(t *testing.T) {
	id := uuid.New()
	entityID := uuid.New()
	eventID := uuid.New()
	instanceID := uuid.New()
	now := time.Now()
	processedAt := now.Add(1 * time.Hour)
	errorMsg := "test error"

	scheduler := domain.Scheduler{
		ID:           id,
		EntityID:     entityID,
		EventID:      eventID,
		InstanceID:   &instanceID,
		Action:       domain.SchedulerActionReminder,
		Status:       domain.SchedulerStatusProcessed,
		ScheduledAt:  now,
		ProcessedAt:  &processedAt,
		Retries:      2,
		MaxRetries:   3,
		ErrorMessage: &errorMsg,
		Metadata:     map[string]interface{}{"type": "email"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	assert.Equal(t, id, scheduler.ID)
	assert.Equal(t, entityID, scheduler.EntityID)
	assert.Equal(t, eventID, scheduler.EventID)
	assert.Equal(t, &instanceID, scheduler.InstanceID)
	assert.Equal(t, domain.SchedulerActionReminder, scheduler.Action)
	assert.Equal(t, domain.SchedulerStatusProcessed, scheduler.Status)
	assert.Equal(t, now, scheduler.ScheduledAt)
	assert.Equal(t, &processedAt, scheduler.ProcessedAt)
	assert.Equal(t, 2, scheduler.Retries)
	assert.Equal(t, 3, scheduler.MaxRetries)
	assert.Equal(t, &errorMsg, scheduler.ErrorMessage)
	assert.Equal(t, "email", scheduler.Metadata["type"])
}

func TestScheduler_NilOptionalFields(t *testing.T) {
	scheduler := domain.Scheduler{
		ID:          uuid.New(),
		EntityID:    uuid.New(),
		EventID:     uuid.New(),
		Action:      domain.SchedulerActionConfirmation,
		Status:      domain.SchedulerStatusPending,
		ScheduledAt: time.Now(),
		MaxRetries:  3,
	}

	assert.Nil(t, scheduler.InstanceID)
	assert.Nil(t, scheduler.ProcessedAt)
	assert.Nil(t, scheduler.ErrorMessage)
	assert.Nil(t, scheduler.Metadata)
}

func TestCreateSchedulerInput(t *testing.T) {
	eventID := uuid.New()
	instanceID := uuid.New()
	scheduledAt := time.Now().Add(24 * time.Hour)

	input := domain.CreateSchedulerInput{
		EventID:     eventID,
		InstanceID:  &instanceID,
		Action:      domain.SchedulerActionReminder,
		ScheduledAt: scheduledAt,
		MaxRetries:  5,
		Metadata:    map[string]interface{}{"channel": "whatsapp"},
	}

	assert.Equal(t, eventID, input.EventID)
	assert.Equal(t, &instanceID, input.InstanceID)
	assert.Equal(t, domain.SchedulerActionReminder, input.Action)
	assert.Equal(t, scheduledAt, input.ScheduledAt)
	assert.Equal(t, 5, input.MaxRetries)
	assert.Equal(t, "whatsapp", input.Metadata["channel"])
}

func TestScheduler_AllActions(t *testing.T) {
	actions := []domain.SchedulerAction{
		domain.SchedulerActionConfirmation,
		domain.SchedulerActionReminder,
		domain.SchedulerActionClosure,
		domain.SchedulerActionLocation,
	}

	for _, action := range actions {
		scheduler := domain.Scheduler{
			ID:          uuid.New(),
			EntityID:    uuid.New(),
			EventID:     uuid.New(),
			Action:      action,
			Status:      domain.SchedulerStatusPending,
			ScheduledAt: time.Now(),
		}
		assert.Equal(t, action, scheduler.Action)
	}
}

func TestScheduler_AllStatuses(t *testing.T) {
	statuses := []domain.SchedulerStatus{
		domain.SchedulerStatusPending,
		domain.SchedulerStatusProcessed,
		domain.SchedulerStatusFailed,
		domain.SchedulerStatusSkipped,
	}

	for _, status := range statuses {
		scheduler := domain.Scheduler{
			ID:          uuid.New(),
			EntityID:    uuid.New(),
			EventID:     uuid.New(),
			Action:      domain.SchedulerActionConfirmation,
			Status:      status,
			ScheduledAt: time.Now(),
		}
		assert.Equal(t, status, scheduler.Status)
	}
}

func TestScheduler_Retries(t *testing.T) {
	tests := []struct {
		name       string
		retries    int
		maxRetries int
	}{
		{
			name:       "No retries yet",
			retries:    0,
			maxRetries: 3,
		},
		{
			name:       "Some retries",
			retries:    2,
			maxRetries: 5,
		},
		{
			name:       "Max retries reached",
			retries:    3,
			maxRetries: 3,
		},
		{
			name:       "No retry limit",
			retries:    10,
			maxRetries: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduler := domain.Scheduler{
				ID:          uuid.New(),
				EntityID:    uuid.New(),
				EventID:     uuid.New(),
				Action:      domain.SchedulerActionReminder,
				Status:      domain.SchedulerStatusPending,
				ScheduledAt: time.Now(),
				Retries:     tt.retries,
				MaxRetries:  tt.maxRetries,
			}

			assert.Equal(t, tt.retries, scheduler.Retries)
			assert.Equal(t, tt.maxRetries, scheduler.MaxRetries)
		})
	}
}
