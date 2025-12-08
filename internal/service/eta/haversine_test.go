package eta_test

import (
	"testing"

	"event-coming/internal/service/eta"

	"github.com/stretchr/testify/assert"
)

func TestCalculateHaversineDistance(t *testing.T) {
	tests := []struct {
		name        string
		lat1, lon1  float64
		lat2, lon2  float64
		expectedMin float64
		expectedMax float64
	}{
		{
			name:        "Same point",
			lat1:        -23.550520,
			lon1:        -46.633308,
			lat2:        -23.550520,
			lon2:        -46.633308,
			expectedMin: 0,
			expectedMax: 0.01,
		},
		{
			name:        "São Paulo to Rio de Janeiro",
			lat1:        -23.550520,
			lon1:        -46.633308,
			lat2:        -22.906847,
			lon2:        -43.172897,
			expectedMin: 350000, // ~357 km
			expectedMax: 365000,
		},
		{
			name:        "New York to Los Angeles",
			lat1:        40.712776,
			lon1:        -74.005974,
			lat2:        34.052235,
			lon2:        -118.243683,
			expectedMin: 3900000, // ~3936 km
			expectedMax: 4000000,
		},
		{
			name:        "London to Paris",
			lat1:        51.507351,
			lon1:        -0.127758,
			lat2:        48.856614,
			lon2:        2.352222,
			expectedMin: 330000, // ~344 km
			expectedMax: 360000,
		},
		{
			name:        "Short distance - 1 km",
			lat1:        -23.550520,
			lon1:        -46.633308,
			lat2:        -23.559520,
			lon2:        -46.633308,
			expectedMin: 900,
			expectedMax: 1100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := eta.CalculateHaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			assert.GreaterOrEqual(t, distance, tt.expectedMin)
			assert.LessOrEqual(t, distance, tt.expectedMax)
		})
	}
}

func TestCalculateHaversineDistance_Symmetry(t *testing.T) {
	lat1, lon1 := -23.550520, -46.633308
	lat2, lon2 := -22.906847, -43.172897

	distanceAtoB := eta.CalculateHaversineDistance(lat1, lon1, lat2, lon2)
	distanceBtoA := eta.CalculateHaversineDistance(lat2, lon2, lat1, lon1)

	// Distance should be the same in both directions
	assert.InDelta(t, distanceAtoB, distanceBtoA, 0.01)
}

func TestCalculateHaversineDistance_PolarCoordinates(t *testing.T) {
	// North Pole to South Pole
	distance := eta.CalculateHaversineDistance(90, 0, -90, 0)

	// Should be approximately half the Earth's circumference (~20,000 km)
	assert.GreaterOrEqual(t, distance, 19000000.0)
	assert.LessOrEqual(t, distance, 21000000.0)
}

func TestCalculateHaversineDistance_CrossEquator(t *testing.T) {
	// Point north of equator to point south of equator
	distance := eta.CalculateHaversineDistance(1.0, 0, -1.0, 0)

	// Should be approximately 222 km (2 degrees latitude)
	assert.GreaterOrEqual(t, distance, 200000.0)
	assert.LessOrEqual(t, distance, 250000.0)
}

func TestCalculateHaversineDistance_DateLine(t *testing.T) {
	// Test crossing the international date line
	distance := eta.CalculateHaversineDistance(0, 179.9, 0, -179.9)

	// Should be approximately 22 km
	assert.GreaterOrEqual(t, distance, 20000.0)
	assert.LessOrEqual(t, distance, 25000.0)
}
