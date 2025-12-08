package dto_test

import (
	"testing"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateParticipantRequest(t *testing.T) {
	email := "participant@example.com"
	instanceID := uuid.New()

	req := dto.CreateParticipantRequest{
		Name:        "Test Participant",
		PhoneNumber: "+5511999999999",
		Email:       &email,
		InstanceID:  &instanceID,
		Metadata:    map[string]interface{}{"vip": true},
	}

	assert.Equal(t, "Test Participant", req.Name)
	assert.Equal(t, "+5511999999999", req.PhoneNumber)
	assert.Equal(t, &email, req.Email)
	assert.Equal(t, &instanceID, req.InstanceID)
	assert.True(t, req.Metadata["vip"].(bool))
}

func TestBatchCreateParticipantsRequest(t *testing.T) {
	email1 := "p1@example.com"
	email2 := "p2@example.com"

	req := dto.BatchCreateParticipantsRequest{
		Participants: []dto.CreateParticipantRequest{
			{
				Name:        "Participant 1",
				PhoneNumber: "+5511999999999",
				Email:       &email1,
			},
			{
				Name:        "Participant 2",
				PhoneNumber: "+5511888888888",
				Email:       &email2,
			},
		},
	}

	assert.Len(t, req.Participants, 2)
	assert.Equal(t, "Participant 1", req.Participants[0].Name)
	assert.Equal(t, "Participant 2", req.Participants[1].Name)
}

func TestUpdateParticipantRequest(t *testing.T) {
	name := "Updated Name"
	phone := "+5511888888888"
	email := "updated@example.com"
	status := domain.ParticipantStatusConfirmed

	req := dto.UpdateParticipantRequest{
		Name:        &name,
		PhoneNumber: &phone,
		Email:       &email,
		Status:      &status,
		Metadata:    map[string]interface{}{"updated": true},
	}

	assert.Equal(t, &name, req.Name)
	assert.Equal(t, &phone, req.PhoneNumber)
	assert.Equal(t, &email, req.Email)
	assert.Equal(t, &status, req.Status)
	assert.True(t, req.Metadata["updated"].(bool))
}

func TestParticipantResponse(t *testing.T) {
	id := uuid.New()
	eventID := uuid.New()
	instanceID := uuid.New()
	entityID := uuid.New()
	now := time.Now()
	confirmedAt := now.Add(-1 * time.Hour)
	checkedInAt := now.Add(-30 * time.Minute)
	email := "participant@example.com"

	resp := dto.ParticipantResponse{
		ID:          id,
		EventID:     eventID,
		InstanceID:  &instanceID,
		EntityID:    entityID,
		Name:        "Test Participant",
		PhoneNumber: "+5511999999999",
		Email:       &email,
		Status:      domain.ParticipantStatusCheckedIn,
		ConfirmedAt: &confirmedAt,
		CheckedInAt: &checkedInAt,
		Metadata:    map[string]interface{}{"vip": true},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, eventID, resp.EventID)
	assert.Equal(t, &instanceID, resp.InstanceID)
	assert.Equal(t, entityID, resp.EntityID)
	assert.Equal(t, "Test Participant", resp.Name)
	assert.Equal(t, "+5511999999999", resp.PhoneNumber)
	assert.Equal(t, &email, resp.Email)
	assert.Equal(t, domain.ParticipantStatusCheckedIn, resp.Status)
	assert.Equal(t, &confirmedAt, resp.ConfirmedAt)
	assert.Equal(t, &checkedInAt, resp.CheckedInAt)
}

func TestToParticipantResponse(t *testing.T) {
	id := uuid.New()
	eventID := uuid.New()
	entityID := uuid.New()
	now := time.Now()

	participant := &domain.Participant{
		ID:        id,
		EventID:   eventID,
		EntityID:  entityID,
		Status:    domain.ParticipantStatusConfirmed,
		Metadata:  map[string]interface{}{"role": "speaker"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := dto.ToParticipantResponse(participant)

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, eventID, resp.EventID)
	assert.Equal(t, entityID, resp.EntityID)
	assert.Equal(t, domain.ParticipantStatusConfirmed, resp.Status)
	assert.Equal(t, "speaker", resp.Metadata["role"])
}

func TestParticipantResponse_NilOptionalFields(t *testing.T) {
	resp := dto.ParticipantResponse{
		ID:          uuid.New(),
		EventID:     uuid.New(),
		EntityID:    uuid.New(),
		Name:        "Minimal Participant",
		PhoneNumber: "+5511999999999",
		Status:      domain.ParticipantStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	assert.Nil(t, resp.InstanceID)
	assert.Nil(t, resp.Email)
	assert.Nil(t, resp.ConfirmedAt)
	assert.Nil(t, resp.CheckedInAt)
	assert.Nil(t, resp.Metadata)
}
