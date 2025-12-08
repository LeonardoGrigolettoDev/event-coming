package domain_test

import (
	"testing"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParticipant_TableName(t *testing.T) {
	participant := domain.Participant{}
	assert.Equal(t, "participants", participant.TableName())
}

func TestParticipantStatus_Constants(t *testing.T) {
	assert.Equal(t, domain.ParticipantStatus("pending"), domain.ParticipantStatusPending)
	assert.Equal(t, domain.ParticipantStatus("confirmed"), domain.ParticipantStatusConfirmed)
	assert.Equal(t, domain.ParticipantStatus("denied"), domain.ParticipantStatusDenied)
	assert.Equal(t, domain.ParticipantStatus("checked_in"), domain.ParticipantStatusCheckedIn)
	assert.Equal(t, domain.ParticipantStatus("no_show"), domain.ParticipantStatusNoShow)
}

func TestParticipant_Fields(t *testing.T) {
	id := uuid.New()
	eventID := uuid.New()
	instanceID := uuid.New()
	entityID := uuid.New()
	refEntityID := uuid.New()
	now := time.Now()
	confirmedAt := now.Add(-1 * time.Hour)
	checkedInAt := now.Add(-30 * time.Minute)

	participant := domain.Participant{
		ID:          id,
		EventID:     eventID,
		InstanceID:  &instanceID,
		EntityID:    entityID,
		RefEntityID: &refEntityID,
		Status:      domain.ParticipantStatusCheckedIn,
		ConfirmedAt: &confirmedAt,
		CheckedInAt: &checkedInAt,
		Metadata:    map[string]interface{}{"vip": true},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, id, participant.ID)
	assert.Equal(t, eventID, participant.EventID)
	assert.Equal(t, &instanceID, participant.InstanceID)
	assert.Equal(t, entityID, participant.EntityID)
	assert.Equal(t, &refEntityID, participant.RefEntityID)
	assert.Equal(t, domain.ParticipantStatusCheckedIn, participant.Status)
	assert.Equal(t, &confirmedAt, participant.ConfirmedAt)
	assert.Equal(t, &checkedInAt, participant.CheckedInAt)
	assert.True(t, participant.Metadata["vip"].(bool))
}

func TestParticipant_NilOptionalFields(t *testing.T) {
	participant := domain.Participant{
		ID:       uuid.New(),
		EventID:  uuid.New(),
		EntityID: uuid.New(),
		Status:   domain.ParticipantStatusPending,
	}

	assert.Nil(t, participant.InstanceID)
	assert.Nil(t, participant.RefEntityID)
	assert.Nil(t, participant.ConfirmedAt)
	assert.Nil(t, participant.CheckedInAt)
	assert.Nil(t, participant.Metadata)
	assert.Nil(t, participant.Entity)
	assert.Nil(t, participant.RefEntity)
}

func TestCreateParticipantInput(t *testing.T) {
	eventID := uuid.New()
	instanceID := uuid.New()
	email := "participant@example.com"

	input := domain.CreateParticipantInput{
		EventID:     eventID,
		InstanceID:  &instanceID,
		Name:        "Test Participant",
		PhoneNumber: "+5511999999999",
		Email:       &email,
		Metadata:    map[string]interface{}{"role": "speaker"},
	}

	assert.Equal(t, eventID, input.EventID)
	assert.Equal(t, &instanceID, input.InstanceID)
	assert.Equal(t, "Test Participant", input.Name)
	assert.Equal(t, "+5511999999999", input.PhoneNumber)
	assert.Equal(t, &email, input.Email)
	assert.Equal(t, "speaker", input.Metadata["role"])
}

func TestUpdateParticipantInput(t *testing.T) {
	name := "Updated Name"
	phone := "+5511888888888"
	email := "updated@example.com"
	status := domain.ParticipantStatusConfirmed

	input := domain.UpdateParticipantInput{
		Name:        &name,
		PhoneNumber: &phone,
		Email:       &email,
		Status:      &status,
		Metadata:    map[string]interface{}{"updated": true},
	}

	assert.Equal(t, &name, input.Name)
	assert.Equal(t, &phone, input.PhoneNumber)
	assert.Equal(t, &email, input.Email)
	assert.Equal(t, &status, input.Status)
	assert.True(t, input.Metadata["updated"].(bool))
}

func TestUpdateParticipantInput_NilFields(t *testing.T) {
	input := domain.UpdateParticipantInput{}

	assert.Nil(t, input.Name)
	assert.Nil(t, input.PhoneNumber)
	assert.Nil(t, input.Email)
	assert.Nil(t, input.Status)
	assert.Nil(t, input.Metadata)
}

func TestParticipantDistance(t *testing.T) {
	id := uuid.New()
	eta := 15
	now := time.Now()

	distance := domain.ParticipantDistance{
		ParticipantID: id,
		Name:          "Nearby Participant",
		Distance:      1500.5,
		ETA:           &eta,
		LastUpdate:    now,
	}

	assert.Equal(t, id, distance.ParticipantID)
	assert.Equal(t, "Nearby Participant", distance.Name)
	assert.Equal(t, 1500.5, distance.Distance)
	assert.Equal(t, &eta, distance.ETA)
	assert.Equal(t, now, distance.LastUpdate)
}

func TestParticipantDistance_NoETA(t *testing.T) {
	distance := domain.ParticipantDistance{
		ParticipantID: uuid.New(),
		Name:          "Participant Without ETA",
		Distance:      5000.0,
		LastUpdate:    time.Now(),
	}

	assert.Nil(t, distance.ETA)
}

func TestParticipant_AllStatuses(t *testing.T) {
	statuses := []domain.ParticipantStatus{
		domain.ParticipantStatusPending,
		domain.ParticipantStatusConfirmed,
		domain.ParticipantStatusDenied,
		domain.ParticipantStatusCheckedIn,
		domain.ParticipantStatusNoShow,
	}

	for _, status := range statuses {
		participant := domain.Participant{
			ID:       uuid.New(),
			EventID:  uuid.New(),
			EntityID: uuid.New(),
			Status:   status,
		}
		assert.Equal(t, status, participant.Status)
	}
}
