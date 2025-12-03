package rrule

import (
	"fmt"
	"strings"
	"time"
)

// Parser handles RRULE parsing and recurrence generation
type Parser struct{}

// NewParser creates a new RRULE parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseRRule parses an RRULE string
// This is a simplified implementation - for production, use github.com/teambition/rrule-go
func (p *Parser) ParseRRule(rrule string) (map[string]string, error) {
	if !strings.HasPrefix(rrule, "RRULE:") {
		return nil, fmt.Errorf("invalid RRULE format: must start with 'RRULE:'")
	}

	rrule = strings.TrimPrefix(rrule, "RRULE:")
	parts := strings.Split(rrule, ";")
	
	result := make(map[string]string)
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}

	return result, nil
}

// GenerateInstances generates event instances based on RRULE
// This is a simplified implementation for common cases
func (p *Parser) GenerateInstances(startTime time.Time, rrule string, until time.Time) ([]time.Time, error) {
	parsed, err := p.ParseRRule(rrule)
	if err != nil {
		return nil, err
	}

	freq, ok := parsed["FREQ"]
	if !ok {
		return nil, fmt.Errorf("FREQ is required in RRULE")
	}

	var instances []time.Time
	current := startTime

	// Simplified generation logic
	switch freq {
	case "DAILY":
		for current.Before(until) {
			instances = append(instances, current)
			current = current.AddDate(0, 0, 1)
		}
	case "WEEKLY":
		for current.Before(until) {
			instances = append(instances, current)
			current = current.AddDate(0, 0, 7)
		}
	case "MONTHLY":
		for current.Before(until) {
			instances = append(instances, current)
			current = current.AddDate(0, 1, 0)
		}
	default:
		return nil, fmt.Errorf("unsupported frequency: %s", freq)
	}

	return instances, nil
}
