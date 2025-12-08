package dto_test

import (
	"testing"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParticipantInput(t *testing.T) {
	email := "participant@example.com"

	input := dto.ParticipantInput{
		Name:        "Test Participant",
		PhoneNumber: "+5511999999999",
		Email:       &email,
		Metadata:    map[string]interface{}{"vip": true},
	}

	assert.Equal(t, "Test Participant", input.Name)
	assert.Equal(t, "+5511999999999", input.PhoneNumber)
	assert.Equal(t, &email, input.Email)
	assert.True(t, input.Metadata["vip"].(bool))
}

func TestSchedulerConfig(t *testing.T) {
	now := time.Now()
	confirmTime := now.Add(-24 * time.Hour)
	reminderTime := now.Add(-2 * time.Hour)
	reminderHours := 2
	locationTime := now.Add(-1 * time.Hour)

	config := dto.SchedulerConfig{
		SendConfirmation:     true,
		ConfirmationTime:     &confirmTime,
		SendReminder:         true,
		ReminderTime:         &reminderTime,
		ReminderBeforeHours:  &reminderHours,
		TrackLocation:        true,
		LocationTrackingTime: &locationTime,
	}

	assert.True(t, config.SendConfirmation)
	assert.Equal(t, &confirmTime, config.ConfirmationTime)
	assert.True(t, config.SendReminder)
	assert.Equal(t, &reminderTime, config.ReminderTime)
	assert.Equal(t, &reminderHours, config.ReminderBeforeHours)
	assert.True(t, config.TrackLocation)
	assert.Equal(t, &locationTime, config.LocationTrackingTime)
}

func TestSchedulerConfig_AllDisabled(t *testing.T) {
	config := dto.SchedulerConfig{
		SendConfirmation: false,
		SendReminder:     false,
		TrackLocation:    false,
	}

	assert.False(t, config.SendConfirmation)
	assert.False(t, config.SendReminder)
	assert.False(t, config.TrackLocation)
	assert.Nil(t, config.ConfirmationTime)
	assert.Nil(t, config.ReminderTime)
	assert.Nil(t, config.ReminderBeforeHours)
	assert.Nil(t, config.LocationTrackingTime)
}

func TestCreateEventRequest(t *testing.T) {
	now := time.Now()
	endTime := now.Add(2 * time.Hour)
	description := "Test event description"
	address := "123 Test Street"
	rrule := "FREQ=WEEKLY;BYDAY=MO"
	deadline := now.Add(-1 * time.Hour)

	email := "p1@example.com"
	participants := []dto.ParticipantInput{
		{
			Name:        "Participant 1",
			PhoneNumber: "+5511999999999",
			Email:       &email,
		},
	}

	scheduler := &dto.SchedulerConfig{
		SendConfirmation: true,
		SendReminder:     true,
	}

	req := dto.CreateEventRequest{
		Name:                 "Test Event",
		Description:          &description,
		Type:                 domain.EventTypePeriodic,
		LocationLat:          -23.550520,
		LocationLng:          -46.633308,
		LocationAddress:      &address,
		StartTime:            now,
		EndTime:              &endTime,
		RRuleString:          &rrule,
		ConfirmationDeadline: &deadline,
		Participants:         participants,
		Scheduler:            scheduler,
	}

	assert.Equal(t, "Test Event", req.Name)
	assert.Equal(t, &description, req.Description)
	assert.Equal(t, domain.EventTypePeriodic, req.Type)
	assert.Equal(t, -23.550520, req.LocationLat)
	assert.Equal(t, -46.633308, req.LocationLng)
	assert.Equal(t, &address, req.LocationAddress)
	assert.Equal(t, now, req.StartTime)
	assert.Equal(t, &endTime, req.EndTime)
	assert.Equal(t, &rrule, req.RRuleString)
	assert.Equal(t, &deadline, req.ConfirmationDeadline)
	assert.Len(t, req.Participants, 1)
	assert.NotNil(t, req.Scheduler)
}

func TestUpdateEventRequest(t *testing.T) {
	now := time.Now()
	endTime := now.Add(3 * time.Hour)
	name := "Updated Event"
	description := "Updated description"
	status := domain.EventStatusActive
	lat := -22.906847
	lng := -43.172897
	address := "456 New Street"
	deadline := now.Add(-15 * time.Minute)

	req := dto.UpdateEventRequest{
		Name:                 &name,
		Description:          &description,
		Status:               &status,
		LocationLat:          &lat,
		LocationLng:          &lng,
		LocationAddress:      &address,
		StartTime:            &now,
		EndTime:              &endTime,
		ConfirmationDeadline: &deadline,
	}

	assert.Equal(t, &name, req.Name)
	assert.Equal(t, &description, req.Description)
	assert.Equal(t, &status, req.Status)
	assert.Equal(t, &lat, req.LocationLat)
	assert.Equal(t, &lng, req.LocationLng)
	assert.Equal(t, &address, req.LocationAddress)
	assert.Equal(t, &now, req.StartTime)
	assert.Equal(t, &endTime, req.EndTime)
	assert.Equal(t, &deadline, req.ConfirmationDeadline)
}

func TestEventResponse(t *testing.T) {
	id := uuid.New()
	entityID := uuid.New()
	createdBy := uuid.New()
	now := time.Now()
	endTime := now.Add(2 * time.Hour)
	description := "Test event"
	address := "123 Test Street"
	rrule := "FREQ=DAILY"
	deadline := now.Add(-1 * time.Hour)

	resp := dto.EventResponse{
		ID:                   id,
		EntityID:             entityID,
		Name:                 "Test Event",
		Description:          &description,
		Type:                 domain.EventTypeDemand,
		Status:               domain.EventStatusScheduled,
		LocationLat:          -23.550520,
		LocationLng:          -46.633308,
		LocationAddress:      &address,
		StartTime:            now,
		EndTime:              &endTime,
		RRuleString:          &rrule,
		ConfirmationDeadline: &deadline,
		CreatedBy:            createdBy,
		CreatedAt:            now,
		UpdatedAt:            now,
		SchedulersCreated:    3,
	}

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, entityID, resp.EntityID)
	assert.Equal(t, "Test Event", resp.Name)
	assert.Equal(t, &description, resp.Description)
	assert.Equal(t, domain.EventTypeDemand, resp.Type)
	assert.Equal(t, domain.EventStatusScheduled, resp.Status)
	assert.Equal(t, -23.550520, resp.LocationLat)
	assert.Equal(t, -46.633308, resp.LocationLng)
	assert.Equal(t, 3, resp.SchedulersCreated)
}

func TestEventResponse_WithParticipants(t *testing.T) {
	id := uuid.New()
	entityID := uuid.New()
	participantID := uuid.New()
	now := time.Now()

	participant := &dto.ParticipantResponse{
		ID:       participantID,
		EventID:  id,
		EntityID: entityID,
		Name:     "Test Participant",
		Status:   domain.ParticipantStatusConfirmed,
	}

	resp := dto.EventResponse{
		ID:           id,
		EntityID:     entityID,
		Name:         "Test Event",
		Type:         domain.EventTypeDemand,
		Status:       domain.EventStatusActive,
		LocationLat:  0,
		LocationLng:  0,
		StartTime:    now,
		CreatedBy:    uuid.New(),
		Participants: []*dto.ParticipantResponse{participant},
	}

	assert.Len(t, resp.Participants, 1)
	assert.Equal(t, participantID, resp.Participants[0].ID)
}

func TestToEventResponse(t *testing.T) {
	id := uuid.New()
	entityID := uuid.New()
	createdBy := uuid.New()
	now := time.Now()
	description := "Test event"

	event := &domain.Event{
		ID:          id,
		EntityID:    entityID,
		Name:        "Test Event",
		Description: &description,
		Type:        domain.EventTypeDemand,
		Status:      domain.EventStatusScheduled,
		LocationLat: -23.550520,
		LocationLng: -46.633308,
		StartTime:   now,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	resp := dto.ToEventResponse(event)

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, entityID, resp.EntityID)
	assert.Equal(t, "Test Event", resp.Name)
	assert.Equal(t, &description, resp.Description)
	assert.Equal(t, domain.EventTypeDemand, resp.Type)
	assert.Equal(t, domain.EventStatusScheduled, resp.Status)
	assert.Equal(t, -23.550520, resp.LocationLat)
	assert.Equal(t, -46.633308, resp.LocationLng)
	assert.Equal(t, createdBy, resp.CreatedBy)
}
