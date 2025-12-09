package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// PubSub gerencia a comunicação entre instâncias via Redis
type PubSub struct {
	client *redis.Client
	hub    *Hub
	logger *zap.Logger
}

// NewPubSub cria um novo gerenciador de PubSub
func NewPubSub(client *redis.Client, hub *Hub, logger *zap.Logger) *PubSub {
	return &PubSub{
		client: client,
		hub:    hub,
		logger: logger,
	}
}

// getRedisChannel retorna o nome do canal Redis para um evento
func getRedisChannel(entityID, eventID string) string {
	return fmt.Sprintf("ws:event:%s:%s", entityID, eventID)
}

// Publish publica uma mensagem no Redis para todas as instâncias
func (p *PubSub) Publish(ctx context.Context, entityID, eventID string, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	channel := getRedisChannel(entityID, eventID)
	if err := p.client.Publish(ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish to Redis: %w", err)
	}

	p.logger.Debug("Published message to Redis",
		zap.String("channel", channel),
		zap.String("type", string(msg.Type)),
	)

	return nil
}

// Subscribe se inscreve em um canal de evento e repassa para o Hub local
func (p *PubSub) Subscribe(ctx context.Context, entityID, eventID string) error {
	channel := getRedisChannel(entityID, eventID)
	pubsub := p.client.Subscribe(ctx, channel)

	// Verificar se a inscrição foi bem-sucedida
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe to channel: %w", err)
	}

	p.logger.Info("Subscribed to Redis channel", zap.String("channel", channel))

	// Processar mensagens em goroutine
	go func() {
		defer pubsub.Close()

		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				p.logger.Info("Unsubscribing from Redis channel", zap.String("channel", channel))
				return

			case redisMsg, ok := <-ch:
				if !ok {
					return
				}

				var msg Message
				if err := json.Unmarshal([]byte(redisMsg.Payload), &msg); err != nil {
					p.logger.Warn("Failed to unmarshal Redis message", zap.Error(err))
					continue
				}

				// Broadcast para clientes locais
				if err := p.hub.Broadcast(entityID, eventID, &msg); err != nil {
					p.logger.Error("Failed to broadcast message", zap.Error(err))
				}
			}
		}
	}()

	return nil
}

// SubscribeAll se inscreve em todos os eventos ativos
// Usa pattern matching do Redis
func (p *PubSub) SubscribeAll(ctx context.Context) error {
	pattern := "ws:event:*"
	pubsub := p.client.PSubscribe(ctx, pattern)

	// Verificar se a inscrição foi bem-sucedida
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe to pattern: %w", err)
	}

	p.logger.Info("Subscribed to Redis pattern", zap.String("pattern", pattern))

	// Processar mensagens em goroutine
	go func() {
		defer pubsub.Close()

		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				p.logger.Info("Unsubscribing from Redis pattern", zap.String("pattern", pattern))
				return

			case redisMsg, ok := <-ch:
				if !ok {
					return
				}

				// Extrair entityID e eventID do canal
				// Formato: ws:event:{entityID}:{eventID}
				var entityID, eventID string
				_, err := fmt.Sscanf(redisMsg.Channel, "ws:event:%s", &entityID)
				if err != nil {
					// Tentar parse manual
					entityID, eventID = parseChannel(redisMsg.Channel)
				}

				var msg Message
				if err := json.Unmarshal([]byte(redisMsg.Payload), &msg); err != nil {
					p.logger.Warn("Failed to unmarshal Redis message", zap.Error(err))
					continue
				}

				// Broadcast para clientes locais
				if err := p.hub.Broadcast(entityID, eventID, &msg); err != nil {
					p.logger.Error("Failed to broadcast message", zap.Error(err))
				}
			}
		}
	}()

	return nil
}

// parseChannel extrai entityID e eventID do nome do canal
func parseChannel(channel string) (entityID, eventID string) {
	// ws:event:{entityID}:{eventID}
	var prefix string
	fmt.Sscanf(channel, "%s:%s:%s", &prefix, &entityID, &eventID)
	return
}

// PublishLocationUpdate publica uma atualização de localização
func (p *PubSub) PublishLocationUpdate(ctx context.Context, entityID, eventID string, data *LocationUpdateData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := &Message{
		Type:      MessageTypeLocationUpdate,
		Timestamp: time.Now(),
		Data:      jsonData,
	}

	return p.Publish(ctx, entityID, eventID, msg)
}
