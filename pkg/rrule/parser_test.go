package rrule_test

import (
	"testing"
	"time"

	"event-coming/pkg/rrule"

	"github.com/stretchr/testify/assert"
)

func TestNewParser(t *testing.T) {
	parser := rrule.NewParser()
	assert.NotNil(t, parser)
}

func TestParser_ParseRRule_Valid(t *testing.T) {
	parser := rrule.NewParser()

	tests := []struct {
		name         string
		rrule        string
		expectedFreq string
		expectedKeys []string
	}{
		{
			name:         "Daily frequency",
			rrule:        "RRULE:FREQ=DAILY",
			expectedFreq: "DAILY",
			expectedKeys: []string{"FREQ"},
		},
		{
			name:         "Weekly frequency",
			rrule:        "RRULE:FREQ=WEEKLY",
			expectedFreq: "WEEKLY",
			expectedKeys: []string{"FREQ"},
		},
		{
			name:         "Monthly frequency",
			rrule:        "RRULE:FREQ=MONTHLY",
			expectedFreq: "MONTHLY",
			expectedKeys: []string{"FREQ"},
		},
		{
			name:         "Weekly with BYDAY",
			rrule:        "RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR",
			expectedFreq: "WEEKLY",
			expectedKeys: []string{"FREQ", "BYDAY"},
		},
		{
			name:         "Daily with COUNT",
			rrule:        "RRULE:FREQ=DAILY;COUNT=10",
			expectedFreq: "DAILY",
			expectedKeys: []string{"FREQ", "COUNT"},
		},
		{
			name:         "Weekly with INTERVAL",
			rrule:        "RRULE:FREQ=WEEKLY;INTERVAL=2",
			expectedFreq: "WEEKLY",
			expectedKeys: []string{"FREQ", "INTERVAL"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.ParseRRule(tt.rrule)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedFreq, result["FREQ"])

			for _, key := range tt.expectedKeys {
				_, exists := result[key]
				assert.True(t, exists, "Key %s should exist", key)
			}
		})
	}
}

func TestParser_ParseRRule_Invalid(t *testing.T) {
	parser := rrule.NewParser()

	tests := []struct {
		name  string
		rrule string
	}{
		{
			name:  "Missing RRULE prefix",
			rrule: "FREQ=DAILY",
		},
		{
			name:  "Empty string",
			rrule: "",
		},
		{
			name:  "Just prefix",
			rrule: "RULE:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.ParseRRule(tt.rrule)
			assert.Error(t, err)
		})
	}
}

func TestParser_GenerateInstances_Daily(t *testing.T) {
	parser := rrule.NewParser()

	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 5, 10, 0, 0, 0, time.UTC)

	instances, err := parser.GenerateInstances(startTime, "RRULE:FREQ=DAILY", until)

	assert.NoError(t, err)
	assert.Len(t, instances, 4) // Jan 1, 2, 3, 4

	// Verify dates
	assert.Equal(t, time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), instances[0])
	assert.Equal(t, time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC), instances[1])
	assert.Equal(t, time.Date(2024, 1, 3, 10, 0, 0, 0, time.UTC), instances[2])
	assert.Equal(t, time.Date(2024, 1, 4, 10, 0, 0, 0, time.UTC), instances[3])
}

func TestParser_GenerateInstances_Weekly(t *testing.T) {
	parser := rrule.NewParser()

	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 29, 10, 0, 0, 0, time.UTC)

	instances, err := parser.GenerateInstances(startTime, "RRULE:FREQ=WEEKLY", until)

	assert.NoError(t, err)
	assert.Len(t, instances, 4) // Jan 1, 8, 15, 22

	// Verify dates
	assert.Equal(t, time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), instances[0])
	assert.Equal(t, time.Date(2024, 1, 8, 10, 0, 0, 0, time.UTC), instances[1])
	assert.Equal(t, time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), instances[2])
	assert.Equal(t, time.Date(2024, 1, 22, 10, 0, 0, 0, time.UTC), instances[3])
}

func TestParser_GenerateInstances_Monthly(t *testing.T) {
	parser := rrule.NewParser()

	startTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	until := time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC)

	instances, err := parser.GenerateInstances(startTime, "RRULE:FREQ=MONTHLY", until)

	assert.NoError(t, err)
	assert.Len(t, instances, 4) // Jan 15, Feb 15, Mar 15, Apr 15

	// Verify months
	assert.Equal(t, time.January, instances[0].Month())
	assert.Equal(t, time.February, instances[1].Month())
	assert.Equal(t, time.March, instances[2].Month())
	assert.Equal(t, time.April, instances[3].Month())
}

func TestParser_GenerateInstances_NoInstances(t *testing.T) {
	parser := rrule.NewParser()

	startTime := time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 5, 10, 0, 0, 0, time.UTC) // Before start

	instances, err := parser.GenerateInstances(startTime, "RRULE:FREQ=DAILY", until)

	assert.NoError(t, err)
	assert.Empty(t, instances)
}

func TestParser_GenerateInstances_InvalidRRule(t *testing.T) {
	parser := rrule.NewParser()

	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC)

	_, err := parser.GenerateInstances(startTime, "INVALID", until)
	assert.Error(t, err)
}

func TestParser_GenerateInstances_MissingFreq(t *testing.T) {
	parser := rrule.NewParser()

	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC)

	_, err := parser.GenerateInstances(startTime, "RRULE:COUNT=10", until)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FREQ is required")
}

func TestParser_GenerateInstances_UnsupportedFreq(t *testing.T) {
	parser := rrule.NewParser()

	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC)

	_, err := parser.GenerateInstances(startTime, "RRULE:FREQ=YEARLY", until)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported frequency")
}

func TestParser_GenerateInstances_PreservesTime(t *testing.T) {
	parser := rrule.NewParser()

	// Specific time: 14:30:45
	startTime := time.Date(2024, 1, 1, 14, 30, 45, 0, time.UTC)
	until := time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)

	instances, err := parser.GenerateInstances(startTime, "RRULE:FREQ=DAILY", until)

	assert.NoError(t, err)
	assert.Len(t, instances, 3)

	// All instances should have the same time
	for _, inst := range instances {
		assert.Equal(t, 14, inst.Hour())
		assert.Equal(t, 30, inst.Minute())
		assert.Equal(t, 45, inst.Second())
	}
}
