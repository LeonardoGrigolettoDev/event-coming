package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/dto"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// EventCacheService gerencia dados em cache do Redis
type EventCacheService struct {
	redisClient *redis.Client
}

// NewEventCacheService cria um novo serviço de cache de eventos
func NewEventCacheService(redisClient *redis.Client) *EventCacheService {
	return &EventCacheService{
		redisClient: redisClient,
	}
}

// GetEventCacheData busca todas as informações em cache de um evento
func (s *EventCacheService) GetEventCacheData(ctx context.Context, entID, eventID uuid.UUID) (*dto.EventCacheResponse, error) {
	data := &dto.EventCacheResponse{
		EntityID:      entID,
		EventID:       eventID,
		Locations:     []dto.ParticipantLocationData{},
		Confirmations: []dto.ParticipantConfirmationData{},
		FetchedAt:     time.Now(),
	}

	// Buscar localizações
	locations, err := s.getLocations(ctx, entID, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}
	data.Locations = locations
	data.TotalLocations = len(locations)

	// Buscar confirmações
	confirmations, err := s.getConfirmations(ctx, entID, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get confirmations: %w", err)
	}
	data.Confirmations = confirmations

	// Contar status
	for _, c := range confirmations {
		switch c.Status {
		case domain.ParticipantStatusConfirmed, domain.ParticipantStatusCheckedIn:
			data.TotalConfirmed++
		case domain.ParticipantStatusPending:
			data.TotalPending++
		case domain.ParticipantStatusDenied:
			data.TotalDenied++
		}
	}

	return data, nil
}

// getLocations busca todas as localizações de participantes de um evento
func (s *EventCacheService) getLocations(ctx context.Context, entID, eventID uuid.UUID) ([]dto.ParticipantLocationData, error) {
	// Pattern: location:latest:{eventID}:*
	pattern := fmt.Sprintf("location:latest:%s:*", eventID)

	var locations []dto.ParticipantLocationData
	var cursor uint64

	for {
		keys, nextCursor, err := s.redisClient.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan keys: %w", err)
		}

		if len(keys) > 0 {
			// Buscar valores em batch
			values, err := s.redisClient.MGet(ctx, keys...).Result()
			if err != nil {
				return nil, fmt.Errorf("failed to get values: %w", err)
			}

			for _, val := range values {
				if val == nil {
					continue
				}

				str, ok := val.(string)
				if !ok {
					continue
				}

				var loc domain.Location
				if err := json.Unmarshal([]byte(str), &loc); err != nil {
					continue
				}

				locations = append(locations, dto.ParticipantLocationData{
					ParticipantID:   loc.ParticipantID,
					ParticipantName: "", // Será preenchido se disponível
					Latitude:        loc.Latitude,
					Longitude:       loc.Longitude,
					Accuracy:        loc.Accuracy,
					Speed:           loc.Speed,
					Heading:         loc.Heading,
					UpdatedAt:       loc.Timestamp,
				})
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return locations, nil
}

// getConfirmations busca todas as confirmações de participantes de um evento
func (s *EventCacheService) getConfirmations(ctx context.Context, entID, eventID uuid.UUID) ([]dto.ParticipantConfirmationData, error) {
	// Pattern: confirmation:{entID}:{eventID}:*
	pattern := fmt.Sprintf("confirmation:%s:%s:*", entID, eventID)

	var confirmations []dto.ParticipantConfirmationData
	var cursor uint64

	for {
		keys, nextCursor, err := s.redisClient.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan keys: %w", err)
		}

		if len(keys) > 0 {
			values, err := s.redisClient.MGet(ctx, keys...).Result()
			if err != nil {
				return nil, fmt.Errorf("failed to get values: %w", err)
			}

			for _, val := range values {
				if val == nil {
					continue
				}

				str, ok := val.(string)
				if !ok {
					continue
				}

				var conf dto.ParticipantConfirmationData
				if err := json.Unmarshal([]byte(str), &conf); err != nil {
					continue
				}

				confirmations = append(confirmations, conf)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return confirmations, nil
}

// SetConfirmation salva uma confirmação no cache
func (s *EventCacheService) SetConfirmation(ctx context.Context, entID, eventID uuid.UUID, participant *domain.Participant) error {
	key := fmt.Sprintf("confirmation:%s:%s:%s", entID, eventID, participant.ID)

	data := dto.ParticipantConfirmationData{
		ParticipantID: participant.ID,
		Status:        participant.Status,
		ConfirmedAt:   participant.ConfirmedAt,
		CheckedInAt:   participant.CheckedInAt,
		UpdatedAt:     time.Now(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal confirmation: %w", err)
	}

	// TTL de 24 horas
	if err := s.redisClient.Set(ctx, key, jsonData, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to set confirmation: %w", err)
	}

	return nil
}

// DeleteConfirmation remove uma confirmação do cache
func (s *EventCacheService) DeleteConfirmation(ctx context.Context, entID, eventID, participantID uuid.UUID) error {
	key := fmt.Sprintf("confirmation:%s:%s:%s", entID, eventID, participantID)
	return s.redisClient.Del(ctx, key).Err()
}

// GetLocationsSummary retorna um resumo rápido das localizações
func (s *EventCacheService) GetLocationsSummary(ctx context.Context, eventID uuid.UUID) (int, error) {
	pattern := fmt.Sprintf("location:latest:%s:*", eventID)

	var count int
	var cursor uint64

	for {
		keys, nextCursor, err := s.redisClient.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return 0, err
		}

		count += len(keys)
		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	return count, nil
}
