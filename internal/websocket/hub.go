package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// Tempo máximo para escrever uma mensagem
	writeWait = 10 * time.Second

	// Tempo máximo para ler o próximo pong
	pongWait = 60 * time.Second

	// Período de envio de pings (deve ser menor que pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Tamanho máximo da mensagem
	maxMessageSize = 4096
)

// MessageType define o tipo de mensagem WebSocket
type MessageType string

const (
	MessageTypeLocationUpdate   MessageType = "location_update"
	MessageTypeETAUpdate        MessageType = "eta_update"
	MessageTypeParticipantJoin  MessageType = "participant_join"
	MessageTypeParticipantLeave MessageType = "participant_leave"
	MessageTypeEventUpdate      MessageType = "event_update"
	MessageTypePing             MessageType = "ping"
	MessageTypePong             MessageType = "pong"
)

// Message representa uma mensagem WebSocket
type Message struct {
	Type      MessageType     `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// LocationUpdateData representa dados de atualização de localização
type LocationUpdateData struct {
	ParticipantID   string   `json:"participant_id"`
	ParticipantName string   `json:"participant_name"`
	Latitude        float64  `json:"latitude"`
	Longitude       float64  `json:"longitude"`
	ETAMinutes      *int     `json:"eta_minutes,omitempty"`
	Distance        *float64 `json:"distance_meters,omitempty"`
}

// Client representa uma conexão WebSocket
type Client struct {
	ID             string
	EntityID string
	EventID        string
	UserID         string
	conn           *websocket.Conn
	send           chan []byte
	hub            *Hub
	logger         *zap.Logger
}

// NewClient cria um novo cliente WebSocket
func NewClient(conn *websocket.Conn, hub *Hub, entityID, eventID, userID string, logger *zap.Logger) *Client {
	return &Client{
		ID:             uuid.New().String(),
		EntityID: entityID,
		EventID:        eventID,
		UserID:         userID,
		conn:           conn,
		send:           make(chan []byte, 256),
		hub:            hub,
		logger:         logger,
	}
}

// ReadPump lê mensagens do WebSocket
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("WebSocket read error", zap.Error(err))
			}
			break
		}

		// Processar mensagem recebida (ping/pong, etc.)
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			c.logger.Warn("Invalid message format", zap.Error(err))
			continue
		}

		// Responder ping com pong
		if msg.Type == MessageTypePing {
			pong := Message{
				Type:      MessageTypePong,
				Timestamp: time.Now(),
			}
			if data, err := json.Marshal(pong); err == nil {
				c.send <- data
			}
		}
	}
}

// WritePump envia mensagens para o WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub fechou o canal
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Hub gerencia todas as conexões WebSocket
type Hub struct {
	// Clientes registrados por evento (org:event -> clients)
	clients    map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *BroadcastMessage
	mu         sync.RWMutex
	logger     *zap.Logger
}

// BroadcastMessage representa uma mensagem para broadcast
type BroadcastMessage struct {
	EntityID string
	EventID        string
	Message        []byte
}

// NewHub cria um novo hub
func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
		logger:     logger,
	}
}

// Run inicia o loop principal do hub
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("Hub stopping")
			return

		case client := <-h.register:
			h.addClient(client)

		case client := <-h.unregister:
			h.removeClient(client)

		case msg := <-h.broadcast:
			h.broadcastToEvent(msg)
		}
	}
}

// getChannelKey retorna a chave do canal para um evento
func getChannelKey(entityID, eventID string) string {
	return entityID + ":" + eventID
}

func (h *Hub) addClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	key := getChannelKey(client.EntityID, client.EventID)
	if h.clients[key] == nil {
		h.clients[key] = make(map[*Client]bool)
	}
	h.clients[key][client] = true

	h.logger.Info("Client connected",
		zap.String("client_id", client.ID),
		zap.String("org_id", client.EntityID),
		zap.String("event_id", client.EventID),
		zap.Int("total_clients", len(h.clients[key])),
	)
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	key := getChannelKey(client.EntityID, client.EventID)
	if clients, ok := h.clients[key]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.send)

			h.logger.Info("Client disconnected",
				zap.String("client_id", client.ID),
				zap.String("org_id", client.EntityID),
				zap.String("event_id", client.EventID),
				zap.Int("remaining_clients", len(clients)),
			)

			// Remove o canal se não há mais clientes
			if len(clients) == 0 {
				delete(h.clients, key)
			}
		}
	}
}

func (h *Hub) broadcastToEvent(msg *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	key := getChannelKey(msg.EntityID, msg.EventID)
	clients, ok := h.clients[key]
	if !ok {
		return
	}

	for client := range clients {
		select {
		case client.send <- msg.Message:
		default:
			// Buffer cheio, fecha a conexão
			close(client.send)
			delete(clients, client)
		}
	}
}

// Broadcast envia uma mensagem para todos os clientes de um evento
func (h *Hub) Broadcast(entityID, eventID string, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.broadcast <- &BroadcastMessage{
		EntityID: entityID,
		EventID:        eventID,
		Message:        data,
	}

	return nil
}

// GetClientCount retorna o número de clientes conectados a um evento
func (h *Hub) GetClientCount(entityID, eventID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	key := getChannelKey(entityID, eventID)
	if clients, ok := h.clients[key]; ok {
		return len(clients)
	}
	return 0
}

// Register registra um cliente
func (h *Hub) Register(client *Client) {
	h.register <- client
}
