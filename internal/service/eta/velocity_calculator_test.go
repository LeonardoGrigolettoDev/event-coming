package eta_test

import (
	"context"
	"testing"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/service/eta"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewVelocityCalculator(t *testing.T) {
	vc := eta.NewVelocityCalculator()
	assert.NotNil(t, vc)
}

func TestVelocityCalculator_CalculateVelocity_TooFewLocations(t *testing.T) {
	vc := eta.NewVelocityCalculator()
	ctx := context.Background()

	// Empty locations
	velocity := vc.CalculateVelocity(ctx, []*domain.Location{})
	assert.Equal(t, 0.0, velocity)

	// Single location
	locations := []*domain.Location{
		{
			ID:        uuid.New(),
			Latitude:  -23.550520,
			Longitude: -46.633308,
			Timestamp: time.Now(),
		},
	}
	velocity = vc.CalculateVelocity(ctx, locations)
	assert.Equal(t, 0.0, velocity)
}

func TestVelocityCalculator_CalculateVelocity_TwoLocations(t *testing.T) {
	vc := eta.NewVelocityCalculator()
	ctx := context.Background()

	now := time.Now()
	// 1 km in 100 seconds = 10 m/s
	locations := []*domain.Location{
		{
			ID:        uuid.New(),
			Latitude:  -23.550520,
			Longitude: -46.633308,
			Timestamp: now,
		},
		{
			ID:        uuid.New(),
			Latitude:  -23.559520, // ~1 km south
			Longitude: -46.633308,
			Timestamp: now.Add(100 * time.Second),
		},
	}

	velocity := vc.CalculateVelocity(ctx, locations)
	// Should be approximately 10 m/s
	assert.GreaterOrEqual(t, velocity, 9.0)
	assert.LessOrEqual(t, velocity, 11.0)
}

func TestVelocityCalculator_CalculateVelocity_MultipleLocations(t *testing.T) {
	vc := eta.NewVelocityCalculator()
	ctx := context.Background()

	now := time.Now()
	locations := []*domain.Location{
		{
			ID:        uuid.New(),
			Latitude:  -23.550520,
			Longitude: -46.633308,
			Timestamp: now,
		},
		{
			ID:        uuid.New(),
			Latitude:  -23.555520, // ~500m
			Longitude: -46.633308,
			Timestamp: now.Add(50 * time.Second),
		},
		{
			ID:        uuid.New(),
			Latitude:  -23.560520, // ~500m more
			Longitude: -46.633308,
			Timestamp: now.Add(100 * time.Second),
		},
	}

	velocity := vc.CalculateVelocity(ctx, locations)
	// Should be approximately 10 m/s
	assert.Greater(t, velocity, 0.0)
}

func TestVelocityCalculator_CalculateVelocity_ZeroTimeDiff(t *testing.T) {
	vc := eta.NewVelocityCalculator()
	ctx := context.Background()

	now := time.Now()
	locations := []*domain.Location{
		{
			ID:        uuid.New(),
			Latitude:  -23.550520,
			Longitude: -46.633308,
			Timestamp: now,
		},
		{
			ID:        uuid.New(),
			Latitude:  -23.559520,
			Longitude: -46.633308,
			Timestamp: now, // Same timestamp
		},
	}

	velocity := vc.CalculateVelocity(ctx, locations)
	assert.Equal(t, 0.0, velocity)
}

func TestVelocityCalculator_CalculateETA(t *testing.T) {
	vc := eta.NewVelocityCalculator()

	tests := []struct {
		name        string
		distance    float64
		velocity    float64
		expectedETA int
	}{
		{
			name:        "Zero velocity",
			distance:    1000,
			velocity:    0,
			expectedETA: 0,
		},
		{
			name:        "1 km at 10 m/s",
			distance:    1000,
			velocity:    10,
			expectedETA: 1, // 100 seconds = 1.67 minutes ~ 1 min
		},
		{
			name:        "10 km at 10 m/s",
			distance:    10000,
			velocity:    10,
			expectedETA: 16, // 1000 seconds = 16.67 minutes
		},
		{
			name:        "5 km at 20 m/s",
			distance:    5000,
			velocity:    20,
			expectedETA: 4, // 250 seconds = 4.17 minutes
		},
		{
			name:        "Short distance with velocity",
			distance:    50,
			velocity:    10,
			expectedETA: 1, // Should return at least 1 minute
		},
		{
			name:        "Long distance",
			distance:    50000,
			velocity:    15,
			expectedETA: 55, // 3333 seconds = 55 minutes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eta := vc.CalculateETA(tt.distance, tt.velocity)
			assert.Equal(t, tt.expectedETA, eta)
		})
	}
}

func TestVelocityCalculator_CalculateETA_MinimumOneMinute(t *testing.T) {
	vc := eta.NewVelocityCalculator()

	// Very short distance with velocity should return at least 1 minute
	result := vc.CalculateETA(10, 10) // 1 second travel time
	assert.GreaterOrEqual(t, result, 1)
}

func TestVelocityCalculator_CalculateVelocity_WalkingSpeed(t *testing.T) {
	vc := eta.NewVelocityCalculator()
	ctx := context.Background()

	now := time.Now()
	// Walking at ~5 km/h (~1.4 m/s): 100m in 72 seconds
	locations := []*domain.Location{
		{
			ID:        uuid.New(),
			Latitude:  0,
			Longitude: 0,
			Timestamp: now,
		},
		{
			ID:        uuid.New(),
			Latitude:  0.0009, // ~100m north
			Longitude: 0,
			Timestamp: now.Add(72 * time.Second),
		},
	}

	velocity := vc.CalculateVelocity(ctx, locations)
	// Should be approximately 1.4 m/s (walking speed)
	assert.GreaterOrEqual(t, velocity, 1.0)
	assert.LessOrEqual(t, velocity, 2.0)
}

func TestVelocityCalculator_CalculateVelocity_DrivingSpeed(t *testing.T) {
	vc := eta.NewVelocityCalculator()
	ctx := context.Background()

	now := time.Now()
	// Driving at ~60 km/h (~16.67 m/s): 1km in 60 seconds
	locations := []*domain.Location{
		{
			ID:        uuid.New(),
			Latitude:  0,
			Longitude: 0,
			Timestamp: now,
		},
		{
			ID:        uuid.New(),
			Latitude:  0.009, // ~1km north
			Longitude: 0,
			Timestamp: now.Add(60 * time.Second),
		},
	}

	velocity := vc.CalculateVelocity(ctx, locations)
	// Should be approximately 16 m/s
	assert.GreaterOrEqual(t, velocity, 14.0)
	assert.LessOrEqual(t, velocity, 20.0)
}
