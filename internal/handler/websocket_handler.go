package handler

import (
	"net/http"

	"event-coming/internal/websocket"

	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implementar validação de origem em produção
		return true
	},
}

// WebSocketHandler gerencia conexões WebSocket
type WebSocketHandler struct {
	hub    *websocket.Hub
	pubsub *websocket.PubSub
	logger *zap.Logger
}

// NewWebSocketHandler cria um novo handler de WebSocket
func NewWebSocketHandler(hub *websocket.Hub, pubsub *websocket.PubSub, logger *zap.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		hub:    hub,
		pubsub: pubsub,
		logger: logger,
	}
}

// HandleConnection processa novas conexões WebSocket
// GET /api/ws/:organization/:event
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	orgID := c.Param("organization")
	eventID := c.Param("event")

	if orgID == "" || eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization and event are required"})
		return
	}

	// Obter userID do contexto (se autenticado)
	userID, _ := c.Get("user_id")
	userIDStr, _ := userID.(string)

	// Upgrade para WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade to WebSocket", zap.Error(err))
		return
	}

	// Criar cliente
	client := websocket.NewClient(conn, h.hub, orgID, eventID, userIDStr, h.logger)

	// Registrar no hub
	h.hub.Register(client)

	// Inscrever no Redis PubSub para este evento (se ainda não inscrito)
	go func() {
		if err := h.pubsub.Subscribe(c.Request.Context(), orgID, eventID); err != nil {
			h.logger.Warn("Failed to subscribe to Redis channel",
				zap.String("org_id", orgID),
				zap.String("event_id", eventID),
				zap.Error(err),
			)
		}
	}()

	// Iniciar goroutines de leitura e escrita
	go client.WritePump()
	go client.ReadPump()

	h.logger.Info("WebSocket connection established",
		zap.String("org_id", orgID),
		zap.String("event_id", eventID),
		zap.String("client_id", client.ID),
	)
}

// GetConnectionCount retorna o número de conexões para um evento
// GET /api/v1/events/:org/:event/connections
func (h *WebSocketHandler) GetConnectionCount(c *gin.Context) {
	orgID := c.Param("organization")
	eventID := c.Param("event")

	count := h.hub.GetClientCount(orgID, eventID)

	c.JSON(http.StatusOK, gin.H{
		"organization_id": orgID,
		"event_id":        eventID,
		"connections":     count,
	})
}
