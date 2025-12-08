package domain_test

import (
	"testing"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLocation_TableName(t *testing.T) {
	location := domain.Location{}
	assert.Equal(t, "locations", location.TableName())
}

func TestLocation_Fields(t *testing.T) {
	id := uuid.New()
	participantID := uuid.New()
	eventID := uuid.New()
	instanceID := uuid.New()
	entityID := uuid.New()
	now := time.Now()
	accuracy := 10.5
	altitude := 800.0
	speed := 5.5
	heading := 180.0

	location := domain.Location{
		ID:            id,
		ParticipantID: participantID,
		EventID:       eventID,
		InstanceID:    &instanceID,
		EntityID:      entityID,
		Latitude:      -23.550520,
		Longitude:     -46.633308,
		Accuracy:      &accuracy,
		Altitude:      &altitude,
		Speed:         &speed,
		Heading:       &heading,
		Timestamp:     now,
		CreatedAt:     now,
	}

	assert.Equal(t, id, location.ID)
	assert.Equal(t, participantID, location.ParticipantID)
	assert.Equal(t, eventID, location.EventID)
	assert.Equal(t, &instanceID, location.InstanceID)
	assert.Equal(t, entityID, location.EntityID)
	assert.Equal(t, -23.550520, location.Latitude)
	assert.Equal(t, -46.633308, location.Longitude)
	assert.Equal(t, &accuracy, location.Accuracy)
	assert.Equal(t, &altitude, location.Altitude)
	assert.Equal(t, &speed, location.Speed)
	assert.Equal(t, &heading, location.Heading)
	assert.Equal(t, now, location.Timestamp)
	assert.Equal(t, now, location.CreatedAt)
}

func TestLocation_NilOptionalFields(t *testing.T) {
	location := domain.Location{
		ID:            uuid.New(),
		ParticipantID: uuid.New(),
		EventID:       uuid.New(),
		EntityID:      uuid.New(),
		Latitude:      -23.550520,
		Longitude:     -46.633308,
		Timestamp:     time.Now(),
	}

	assert.Nil(t, location.InstanceID)
	assert.Nil(t, location.Accuracy)
	assert.Nil(t, location.Altitude)
	assert.Nil(t, location.Speed)
	assert.Nil(t, location.Heading)
}

func TestCreateLocationInput(t *testing.T) {
	participantID := uuid.New()
	eventID := uuid.New()
	instanceID := uuid.New()
	now := time.Now()
	accuracy := 5.0
	altitude := 900.0
	speed := 10.0
	heading := 90.0

	input := domain.CreateLocationInput{
		ParticipantID: participantID,
		EventID:       eventID,
		InstanceID:    &instanceID,
		Latitude:      -22.906847,
		Longitude:     -43.172897,
		Accuracy:      &accuracy,
		Altitude:      &altitude,
		Speed:         &speed,
		Heading:       &heading,
		Timestamp:     &now,
	}

	assert.Equal(t, participantID, input.ParticipantID)
	assert.Equal(t, eventID, input.EventID)
	assert.Equal(t, &instanceID, input.InstanceID)
	assert.Equal(t, -22.906847, input.Latitude)
	assert.Equal(t, -43.172897, input.Longitude)
	assert.Equal(t, &accuracy, input.Accuracy)
	assert.Equal(t, &altitude, input.Altitude)
	assert.Equal(t, &speed, input.Speed)
	assert.Equal(t, &heading, input.Heading)
	assert.Equal(t, &now, input.Timestamp)
}

func TestCreateLocationInput_MinimalFields(t *testing.T) {
	input := domain.CreateLocationInput{
		ParticipantID: uuid.New(),
		EventID:       uuid.New(),
		Latitude:      -23.550520,
		Longitude:     -46.633308,
	}

	assert.NotEqual(t, uuid.Nil, input.ParticipantID)
	assert.NotEqual(t, uuid.Nil, input.EventID)
	assert.Equal(t, -23.550520, input.Latitude)
	assert.Equal(t, -46.633308, input.Longitude)
	assert.Nil(t, input.InstanceID)
	assert.Nil(t, input.Accuracy)
	assert.Nil(t, input.Altitude)
	assert.Nil(t, input.Speed)
	assert.Nil(t, input.Heading)
	assert.Nil(t, input.Timestamp)
}

func TestLocation_Coordinates(t *testing.T) {
	tests := []struct {
		name string
		lat  float64
		lng  float64
	}{
		{
			name: "São Paulo",
			lat:  -23.550520,
			lng:  -46.633308,
		},
		{
			name: "Rio de Janeiro",
			lat:  -22.906847,
			lng:  -43.172897,
		},
		{
			name: "New York",
			lat:  40.712776,
			lng:  -74.005974,
		},
		{
			name: "Tokyo",
			lat:  35.689487,
			lng:  139.691711,
		},
		{
			name: "Sydney",
			lat:  -33.868820,
			lng:  151.209290,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := domain.Location{
				ID:            uuid.New(),
				ParticipantID: uuid.New(),
				EventID:       uuid.New(),
				EntityID:      uuid.New(),
				Latitude:      tt.lat,
				Longitude:     tt.lng,
				Timestamp:     time.Now(),
			}

			assert.Equal(t, tt.lat, location.Latitude)
			assert.Equal(t, tt.lng, location.Longitude)
		})
	}
}

func TestLocation_Speed_Heading_Ranges(t *testing.T) {
	tests := []struct {
		name    string
		speed   float64
		heading float64
	}{
		{
			name:    "Stationary",
			speed:   0,
			heading: 0,
		},
		{
			name:    "Walking North",
			speed:   1.4,
			heading: 0,
		},
		{
			name:    "Running East",
			speed:   3.0,
			heading: 90,
		},
		{
			name:    "Cycling South",
			speed:   5.0,
			heading: 180,
		},
		{
			name:    "Driving West",
			speed:   16.67,
			heading: 270,
		},
		{
			name:    "Max Heading",
			speed:   0,
			heading: 360,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := domain.Location{
				ID:            uuid.New(),
				ParticipantID: uuid.New(),
				EventID:       uuid.New(),
				EntityID:      uuid.New(),
				Latitude:      0,
				Longitude:     0,
				Speed:         &tt.speed,
				Heading:       &tt.heading,
				Timestamp:     time.Now(),
			}

			assert.Equal(t, tt.speed, *location.Speed)
			assert.Equal(t, tt.heading, *location.Heading)
		})
	}
}
