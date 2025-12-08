package dto_test

import (
	"testing"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/dto"
	"event-coming/internal/service/eta"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateLocationRequest(t *testing.T) {
	now := time.Now()
	accuracy := 10.5
	altitude := 800.0
	speed := 5.5
	heading := 180.0

	req := dto.CreateLocationRequest{
		Latitude:  -23.550520,
		Longitude: -46.633308,
		Accuracy:  &accuracy,
		Altitude:  &altitude,
		Speed:     &speed,
		Heading:   &heading,
		Timestamp: &now,
	}

	assert.Equal(t, -23.550520, req.Latitude)
	assert.Equal(t, -46.633308, req.Longitude)
	assert.Equal(t, &accuracy, req.Accuracy)
	assert.Equal(t, &altitude, req.Altitude)
	assert.Equal(t, &speed, req.Speed)
	assert.Equal(t, &heading, req.Heading)
	assert.Equal(t, &now, req.Timestamp)
}

func TestCreateLocationRequest_MinimalFields(t *testing.T) {
	req := dto.CreateLocationRequest{
		Latitude:  -23.550520,
		Longitude: -46.633308,
	}

	assert.Equal(t, -23.550520, req.Latitude)
	assert.Equal(t, -46.633308, req.Longitude)
	assert.Nil(t, req.Accuracy)
	assert.Nil(t, req.Altitude)
	assert.Nil(t, req.Speed)
	assert.Nil(t, req.Heading)
	assert.Nil(t, req.Timestamp)
}

func TestLocationResponse(t *testing.T) {
	id := uuid.New()
	participantID := uuid.New()
	eventID := uuid.New()
	entityID := uuid.New()
	now := time.Now()
	accuracy := 5.0
	altitude := 900.0
	speed := 10.0
	heading := 90.0

	resp := dto.LocationResponse{
		ID:            id,
		ParticipantID: participantID,
		EventID:       eventID,
		EntityID:      entityID,
		Latitude:      -22.906847,
		Longitude:     -43.172897,
		Accuracy:      &accuracy,
		Altitude:      &altitude,
		Speed:         &speed,
		Heading:       &heading,
		Timestamp:     now,
		CreatedAt:     now,
	}

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, participantID, resp.ParticipantID)
	assert.Equal(t, eventID, resp.EventID)
	assert.Equal(t, entityID, resp.EntityID)
	assert.Equal(t, -22.906847, resp.Latitude)
	assert.Equal(t, -43.172897, resp.Longitude)
	assert.Equal(t, &accuracy, resp.Accuracy)
	assert.Equal(t, &altitude, resp.Altitude)
	assert.Equal(t, &speed, resp.Speed)
	assert.Equal(t, &heading, resp.Heading)
}

func TestToLocationResponse(t *testing.T) {
	id := uuid.New()
	participantID := uuid.New()
	eventID := uuid.New()
	entityID := uuid.New()
	now := time.Now()

	location := &domain.Location{
		ID:            id,
		ParticipantID: participantID,
		EventID:       eventID,
		EntityID:      entityID,
		Latitude:      -23.550520,
		Longitude:     -46.633308,
		Timestamp:     now,
		CreatedAt:     now,
	}

	resp := dto.ToLocationResponse(location)

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, participantID, resp.ParticipantID)
	assert.Equal(t, eventID, resp.EventID)
	assert.Equal(t, entityID, resp.EntityID)
	assert.Equal(t, -23.550520, resp.Latitude)
	assert.Equal(t, -46.633308, resp.Longitude)
}

func TestToLocationResponse_Nil(t *testing.T) {
	resp := dto.ToLocationResponse(nil)
	assert.Nil(t, resp)
}

func TestToLocationResponseList(t *testing.T) {
	now := time.Now()

	locations := []*domain.Location{
		{
			ID:            uuid.New(),
			ParticipantID: uuid.New(),
			EventID:       uuid.New(),
			EntityID:      uuid.New(),
			Latitude:      -23.550520,
			Longitude:     -46.633308,
			Timestamp:     now,
			CreatedAt:     now,
		},
		{
			ID:            uuid.New(),
			ParticipantID: uuid.New(),
			EventID:       uuid.New(),
			EntityID:      uuid.New(),
			Latitude:      -22.906847,
			Longitude:     -43.172897,
			Timestamp:     now,
			CreatedAt:     now,
		},
	}

	responses := dto.ToLocationResponseList(locations)

	assert.Len(t, responses, 2)
	assert.Equal(t, -23.550520, responses[0].Latitude)
	assert.Equal(t, -22.906847, responses[1].Latitude)
}

func TestToLocationResponseList_Empty(t *testing.T) {
	locations := []*domain.Location{}
	responses := dto.ToLocationResponseList(locations)
	assert.Empty(t, responses)
}

func TestETAResponse(t *testing.T) {
	participantID := uuid.New()
	now := time.Now()

	resp := dto.ETAResponse{
		ParticipantID:  participantID,
		DistanceMeters: 5000.0,
		ETAMinutes:     15,
		Method:         "velocity",
		LastUpdate:     now,
	}

	assert.Equal(t, participantID, resp.ParticipantID)
	assert.Equal(t, 5000.0, resp.DistanceMeters)
	assert.Equal(t, 15, resp.ETAMinutes)
	assert.Equal(t, "velocity", resp.Method)
	assert.Equal(t, now, resp.LastUpdate)
}

func TestToETAResponse(t *testing.T) {
	participantID := uuid.New()
	now := time.Now()

	result := &eta.ETAResult{
		ParticipantID:  participantID,
		DistanceMeters: 3000.0,
		ETAMinutes:     10,
		Method:         "haversine",
		LastUpdate:     now,
	}

	resp := dto.ToETAResponse(result)

	assert.Equal(t, participantID, resp.ParticipantID)
	assert.Equal(t, 3000.0, resp.DistanceMeters)
	assert.Equal(t, 10, resp.ETAMinutes)
	assert.Equal(t, "haversine", resp.Method)
}

func TestToETAResponse_Nil(t *testing.T) {
	resp := dto.ToETAResponse(nil)
	assert.Nil(t, resp)
}

func TestToETAResponseList(t *testing.T) {
	now := time.Now()

	results := []*eta.ETAResult{
		{
			ParticipantID:  uuid.New(),
			DistanceMeters: 2000.0,
			ETAMinutes:     5,
			Method:         "velocity",
			LastUpdate:     now,
		},
		{
			ParticipantID:  uuid.New(),
			DistanceMeters: 8000.0,
			ETAMinutes:     25,
			Method:         "haversine",
			LastUpdate:     now,
		},
	}

	responses := dto.ToETAResponseList(results)

	assert.Len(t, responses, 2)
	assert.Equal(t, 2000.0, responses[0].DistanceMeters)
	assert.Equal(t, 8000.0, responses[1].DistanceMeters)
}

func TestEventETAResponse(t *testing.T) {
	eventID := uuid.New()
	entityID := uuid.New()
	now := time.Now()

	etaResp := &dto.ETAResponse{
		ParticipantID:  uuid.New(),
		DistanceMeters: 1500.0,
		ETAMinutes:     5,
		Method:         "velocity",
		LastUpdate:     now,
	}

	resp := dto.EventETAResponse{
		EventID:      eventID,
		EntityID:     entityID,
		Participants: []*dto.ETAResponse{etaResp},
		FetchedAt:    now,
	}

	assert.Equal(t, eventID, resp.EventID)
	assert.Equal(t, entityID, resp.EntityID)
	assert.Len(t, resp.Participants, 1)
	assert.Equal(t, now, resp.FetchedAt)
}

func TestToEventETAResponse(t *testing.T) {
	eventID := uuid.New()
	entityID := uuid.New()
	now := time.Now()

	results := []*eta.ETAResult{
		{
			ParticipantID:  uuid.New(),
			DistanceMeters: 1000.0,
			ETAMinutes:     3,
			Method:         "velocity",
			LastUpdate:     now,
		},
	}

	resp := dto.ToEventETAResponse(eventID, entityID, results)

	assert.Equal(t, eventID, resp.EventID)
	assert.Equal(t, entityID, resp.EntityID)
	assert.Len(t, resp.Participants, 1)
	assert.NotZero(t, resp.FetchedAt)
}
