package dto_test

import (
	"testing"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParticipantLocationData(t *testing.T) {
	participantID := uuid.New()
	now := time.Now()
	accuracy := 10.5
	speed := 5.5
	heading := 180.0
	eta := now.Add(15 * time.Minute)
	etaMinutes := 15

	data := dto.ParticipantLocationData{
		ParticipantID:   participantID,
		ParticipantName: "Test Participant",
		Latitude:        -23.550520,
		Longitude:       -46.633308,
		Accuracy:        &accuracy,
		Speed:           &speed,
		Heading:         &heading,
		UpdatedAt:       now,
		ETA:             &eta,
		ETAMinutes:      &etaMinutes,
	}

	assert.Equal(t, participantID, data.ParticipantID)
	assert.Equal(t, "Test Participant", data.ParticipantName)
	assert.Equal(t, -23.550520, data.Latitude)
	assert.Equal(t, -46.633308, data.Longitude)
	assert.Equal(t, &accuracy, data.Accuracy)
	assert.Equal(t, &speed, data.Speed)
	assert.Equal(t, &heading, data.Heading)
	assert.Equal(t, now, data.UpdatedAt)
	assert.Equal(t, &eta, data.ETA)
	assert.Equal(t, &etaMinutes, data.ETAMinutes)
}

func TestParticipantLocationData_NilOptionalFields(t *testing.T) {
	data := dto.ParticipantLocationData{
		ParticipantID:   uuid.New(),
		ParticipantName: "Minimal Participant",
		Latitude:        0,
		Longitude:       0,
		UpdatedAt:       time.Now(),
	}

	assert.Nil(t, data.Accuracy)
	assert.Nil(t, data.Speed)
	assert.Nil(t, data.Heading)
	assert.Nil(t, data.ETA)
	assert.Nil(t, data.ETAMinutes)
}

func TestParticipantConfirmationData(t *testing.T) {
	participantID := uuid.New()
	now := time.Now()
	confirmedAt := now.Add(-1 * time.Hour)
	checkedInAt := now.Add(-30 * time.Minute)

	data := dto.ParticipantConfirmationData{
		ParticipantID:   participantID,
		ParticipantName: "Test Participant",
		PhoneNumber:     "+5511999999999",
		Status:          domain.ParticipantStatusCheckedIn,
		ConfirmedAt:     &confirmedAt,
		CheckedInAt:     &checkedInAt,
		UpdatedAt:       now,
	}

	assert.Equal(t, participantID, data.ParticipantID)
	assert.Equal(t, "Test Participant", data.ParticipantName)
	assert.Equal(t, "+5511999999999", data.PhoneNumber)
	assert.Equal(t, domain.ParticipantStatusCheckedIn, data.Status)
	assert.Equal(t, &confirmedAt, data.ConfirmedAt)
	assert.Equal(t, &checkedInAt, data.CheckedInAt)
	assert.Equal(t, now, data.UpdatedAt)
}

func TestParticipantConfirmationData_AllStatuses(t *testing.T) {
	statuses := []domain.ParticipantStatus{
		domain.ParticipantStatusPending,
		domain.ParticipantStatusConfirmed,
		domain.ParticipantStatusDenied,
		domain.ParticipantStatusCheckedIn,
		domain.ParticipantStatusNoShow,
	}

	for _, status := range statuses {
		data := dto.ParticipantConfirmationData{
			ParticipantID:   uuid.New(),
			ParticipantName: "Participant",
			PhoneNumber:     "+5511999999999",
			Status:          status,
			UpdatedAt:       time.Now(),
		}
		assert.Equal(t, status, data.Status)
	}
}

func TestEventCacheResponse(t *testing.T) {
	entityID := uuid.New()
	eventID := uuid.New()
	now := time.Now()

	locations := []dto.ParticipantLocationData{
		{
			ParticipantID:   uuid.New(),
			ParticipantName: "Participant 1",
			Latitude:        -23.550520,
			Longitude:       -46.633308,
			UpdatedAt:       now,
		},
		{
			ParticipantID:   uuid.New(),
			ParticipantName: "Participant 2",
			Latitude:        -22.906847,
			Longitude:       -43.172897,
			UpdatedAt:       now,
		},
	}

	confirmations := []dto.ParticipantConfirmationData{
		{
			ParticipantID:   uuid.New(),
			ParticipantName: "Confirmed 1",
			PhoneNumber:     "+5511999999999",
			Status:          domain.ParticipantStatusConfirmed,
			UpdatedAt:       now,
		},
	}

	resp := dto.EventCacheResponse{
		EntityID:       entityID,
		EventID:        eventID,
		Locations:      locations,
		Confirmations:  confirmations,
		TotalLocations: 2,
		TotalConfirmed: 5,
		TotalPending:   3,
		TotalDenied:    1,
		FetchedAt:      now,
	}

	assert.Equal(t, entityID, resp.EntityID)
	assert.Equal(t, eventID, resp.EventID)
	assert.Len(t, resp.Locations, 2)
	assert.Len(t, resp.Confirmations, 1)
	assert.Equal(t, 2, resp.TotalLocations)
	assert.Equal(t, 5, resp.TotalConfirmed)
	assert.Equal(t, 3, resp.TotalPending)
	assert.Equal(t, 1, resp.TotalDenied)
	assert.Equal(t, now, resp.FetchedAt)
}

func TestEventCacheResponse_Empty(t *testing.T) {
	resp := dto.EventCacheResponse{
		EntityID:       uuid.New(),
		EventID:        uuid.New(),
		Locations:      []dto.ParticipantLocationData{},
		Confirmations:  []dto.ParticipantConfirmationData{},
		TotalLocations: 0,
		TotalConfirmed: 0,
		TotalPending:   0,
		TotalDenied:    0,
		FetchedAt:      time.Now(),
	}

	assert.Empty(t, resp.Locations)
	assert.Empty(t, resp.Confirmations)
	assert.Zero(t, resp.TotalLocations)
	assert.Zero(t, resp.TotalConfirmed)
	assert.Zero(t, resp.TotalPending)
	assert.Zero(t, resp.TotalDenied)
}
