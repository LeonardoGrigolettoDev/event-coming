package eta

import (
	"context"
	"time"

	"event-coming/internal/domain"
)

// VelocityCalculator calculates velocity from location history
type VelocityCalculator struct{}

// NewVelocityCalculator creates a new velocity calculator
func NewVelocityCalculator() *VelocityCalculator {
	return &VelocityCalculator{}
}

// CalculateVelocity calculates average velocity from recent locations
// Returns velocity in meters per second
func (vc *VelocityCalculator) CalculateVelocity(ctx context.Context, locations []*domain.Location) float64 {
	if len(locations) < 2 {
		return 0
	}

	var totalDistance float64
	var totalTime float64

	for i := 1; i < len(locations); i++ {
		prev := locations[i-1]
		curr := locations[i]

		distance := CalculateHaversineDistance(
			prev.Latitude, prev.Longitude,
			curr.Latitude, curr.Longitude,
		)

		timeDiff := curr.Timestamp.Sub(prev.Timestamp).Seconds()
		if timeDiff > 0 {
			totalDistance += distance
			totalTime += timeDiff
		}
	}

	if totalTime == 0 {
		return 0
	}

	return totalDistance / totalTime
}

// CalculateETA calculates estimated time of arrival
// Returns ETA in minutes
func (vc *VelocityCalculator) CalculateETA(distance, velocity float64) int {
	if velocity == 0 {
		return 0
	}

	// Calculate time in seconds
	timeSeconds := distance / velocity

	// Convert to minutes and round up
	timeMinutes := int(timeSeconds / 60)
	if timeMinutes == 0 && distance > 0 {
		timeMinutes = 1
	}

	return timeMinutes
}
