package domain_test

import (
	"testing"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEvent_TableName(t *testing.T) {
	event := domain.Event{}
	assert.Equal(t, "events", event.TableName())
}

func TestEventInstance_TableName(t *testing.T) {
	instance := domain.EventInstance{}
	assert.Equal(t, "event_instances", instance.TableName())
}

func TestEventType_Constants(t *testing.T) {
	assert.Equal(t, domain.EventType("demand"), domain.EventTypeDemand)
	assert.Equal(t, domain.EventType("periodic"), domain.EventTypePeriodic)
}

func TestEventStatus_Constants(t *testing.T) {
	assert.Equal(t, domain.EventStatus("draft"), domain.EventStatusDraft)
	assert.Equal(t, domain.EventStatus("scheduled"), domain.EventStatusScheduled)
	assert.Equal(t, domain.EventStatus("active"), domain.EventStatusActive)
	assert.Equal(t, domain.EventStatus("completed"), domain.EventStatusCompleted)
	assert.Equal(t, domain.EventStatus("cancelled"), domain.EventStatusCancelled)
}

func TestEvent_Fields(t *testing.T) {
	id := uuid.New()
	entityID := uuid.New()
	createdBy := uuid.New()
	now := time.Now()
	description := "Test event description"
	address := "123 Test Street"
	rrule := "FREQ=WEEKLY;BYDAY=MO,WE,FR"
	endTime := now.Add(2 * time.Hour)
	deadline := now.Add(-1 * time.Hour)

	event := domain.Event{
		ID:                   id,
		EntityID:             entityID,
		Name:                 "Test Event",
		Description:          &description,
		Type:                 domain.EventTypePeriodic,
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
	}

	assert.Equal(t, id, event.ID)
	assert.Equal(t, entityID, event.EntityID)
	assert.Equal(t, "Test Event", event.Name)
	assert.Equal(t, &description, event.Description)
	assert.Equal(t, domain.EventTypePeriodic, event.Type)
	assert.Equal(t, domain.EventStatusScheduled, event.Status)
	assert.Equal(t, -23.550520, event.LocationLat)
	assert.Equal(t, -46.633308, event.LocationLng)
	assert.Equal(t, &address, event.LocationAddress)
	assert.Equal(t, now, event.StartTime)
	assert.Equal(t, &endTime, event.EndTime)
	assert.Equal(t, &rrule, event.RRuleString)
	assert.Equal(t, &deadline, event.ConfirmationDeadline)
	assert.Equal(t, createdBy, event.CreatedBy)
}

func TestEventInstance_Fields(t *testing.T) {
	id := uuid.New()
	eventID := uuid.New()
	entityID := uuid.New()
	now := time.Now()
	endTime := now.Add(2 * time.Hour)

	instance := domain.EventInstance{
		ID:           id,
		EventID:      eventID,
		EntityID:     entityID,
		InstanceDate: now,
		Status:       domain.EventStatusActive,
		StartTime:    now,
		EndTime:      &endTime,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	assert.Equal(t, id, instance.ID)
	assert.Equal(t, eventID, instance.EventID)
	assert.Equal(t, entityID, instance.EntityID)
	assert.Equal(t, now, instance.InstanceDate)
	assert.Equal(t, domain.EventStatusActive, instance.Status)
	assert.Equal(t, now, instance.StartTime)
	assert.Equal(t, &endTime, instance.EndTime)
}

func TestCreateEventInput(t *testing.T) {
	now := time.Now()
	endTime := now.Add(2 * time.Hour)
	description := "Test description"
	address := "123 Test Street"
	rrule := "FREQ=DAILY"
	deadline := now.Add(-30 * time.Minute)

	input := domain.CreateEventInput{
		Name:                 "New Event",
		Description:          &description,
		Type:                 domain.EventTypeDemand,
		LocationLat:          -23.550520,
		LocationLng:          -46.633308,
		LocationAddress:      &address,
		StartTime:            now,
		EndTime:              &endTime,
		RRuleString:          &rrule,
		ConfirmationDeadline: &deadline,
	}

	assert.Equal(t, "New Event", input.Name)
	assert.Equal(t, &description, input.Description)
	assert.Equal(t, domain.EventTypeDemand, input.Type)
	assert.Equal(t, -23.550520, input.LocationLat)
	assert.Equal(t, -46.633308, input.LocationLng)
	assert.Equal(t, &address, input.LocationAddress)
	assert.Equal(t, now, input.StartTime)
	assert.Equal(t, &endTime, input.EndTime)
	assert.Equal(t, &rrule, input.RRuleString)
	assert.Equal(t, &deadline, input.ConfirmationDeadline)
}

func TestUpdateEventInput(t *testing.T) {
	now := time.Now()
	endTime := now.Add(3 * time.Hour)
	name := "Updated Event"
	description := "Updated description"
	status := domain.EventStatusActive
	lat := -22.906847
	lng := -43.172897
	address := "456 New Street"
	deadline := now.Add(-15 * time.Minute)

	input := domain.UpdateEventInput{
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

	assert.Equal(t, &name, input.Name)
	assert.Equal(t, &description, input.Description)
	assert.Equal(t, &status, input.Status)
	assert.Equal(t, &lat, input.LocationLat)
	assert.Equal(t, &lng, input.LocationLng)
	assert.Equal(t, &address, input.LocationAddress)
	assert.Equal(t, &now, input.StartTime)
	assert.Equal(t, &endTime, input.EndTime)
	assert.Equal(t, &deadline, input.ConfirmationDeadline)
}

func TestEvent_NilOptionalFields(t *testing.T) {
	event := domain.Event{
		ID:          uuid.New(),
		EntityID:    uuid.New(),
		Name:        "Minimal Event",
		Type:        domain.EventTypeDemand,
		Status:      domain.EventStatusDraft,
		LocationLat: 0,
		LocationLng: 0,
		StartTime:   time.Now(),
		CreatedBy:   uuid.New(),
	}

	assert.Nil(t, event.Description)
	assert.Nil(t, event.LocationAddress)
	assert.Nil(t, event.EndTime)
	assert.Nil(t, event.RRuleString)
	assert.Nil(t, event.ConfirmationDeadline)
	assert.Nil(t, event.Entity)
}

func TestEventInstance_NilOptionalFields(t *testing.T) {
	instance := domain.EventInstance{
		ID:           uuid.New(),
		EventID:      uuid.New(),
		EntityID:     uuid.New(),
		InstanceDate: time.Now(),
		Status:       domain.EventStatusScheduled,
		StartTime:    time.Now(),
	}

	assert.Nil(t, instance.EndTime)
}
