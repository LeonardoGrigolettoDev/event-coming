package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"time"

	"event-coming/internal/config"
	"event-coming/internal/domain"
	"event-coming/internal/dto"
	"event-coming/internal/service"
	"event-coming/internal/whatsapp"
	"event-coming/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// WebhookHandler handles WhatsApp webhook requests
type WebhookHandler struct {
	cfg                *config.WhatsAppConfig
	participantService *service.ParticipantService
	locationService    *service.LocationService
	logger             *zap.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(
	cfg *config.WhatsAppConfig,
	participantService *service.ParticipantService,
	locationService *service.LocationService,
	logger *zap.Logger,
) *WebhookHandler {
	return &WebhookHandler{
		cfg:                cfg,
		participantService: participantService,
		locationService:    locationService,
		logger:             logger,
	}
}

// VerifyWebhook handles webhook verification from WhatsApp
// GET /webhook/whatsapp
func (h *WebhookHandler) VerifyWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == h.cfg.WebhookVerifyToken {
		h.logger.Info("Webhook verified successfully")
		c.String(http.StatusOK, challenge)
		return
	}

	h.logger.Warn("Webhook verification failed",
		zap.String("mode", mode),
		zap.String("token", token),
	)
	response.Error(c, http.StatusForbidden, "forbidden", "Verification failed")
}

// HandleWebhook processes incoming WhatsApp webhook events
// POST /webhook/whatsapp
func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	// Read body for signature verification
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error("Failed to read webhook body", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "bad_request", "Failed to read body")
		return
	}

	// Verify signature if webhook secret is configured
	if h.cfg.WebhookSecret != "" {
		signature := c.GetHeader("X-Hub-Signature-256")
		if !h.verifySignature(body, signature) {
			h.logger.Warn("Invalid webhook signature")
			response.Error(c, http.StatusUnauthorized, "unauthorized", "Invalid signature")
			return
		}
	}

	// Parse payload
	var payload whatsapp.WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Error("Failed to parse webhook payload", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid payload")
		return
	}

	// Process messages
	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if change.Field == "messages" {
				h.processMessages(c, change.Value)
			}
		}
	}

	// Always return 200 to acknowledge receipt
	c.Status(http.StatusOK)
}

// processMessages processes incoming messages
func (h *WebhookHandler) processMessages(c *gin.Context, value whatsapp.Value) {
	for _, msg := range value.Messages {
		switch msg.Type {
		case "location":
			h.handleLocationMessage(c, msg)
		case "interactive":
			h.handleInteractiveMessage(c, msg)
		case "button":
			h.handleButtonMessage(c, msg)
		case "text":
			h.handleTextMessage(c, msg)
		}
	}
}

// handleLocationMessage processes location messages from participants
func (h *WebhookHandler) handleLocationMessage(c *gin.Context, msg whatsapp.Message) {
	if msg.Location == nil {
		return
	}

	phoneNumber := msg.From
	h.logger.Info("Received location from WhatsApp",
		zap.String("phone", phoneNumber),
		zap.Float64("lat", msg.Location.Latitude),
		zap.Float64("lng", msg.Location.Longitude),
	)

	// Find participant by phone number
	participant, err := h.participantService.GetByPhoneNumber(c.Request.Context(), phoneNumber)
	if err != nil {
		h.logger.Warn("Participant not found for phone number",
			zap.String("phone", phoneNumber),
			zap.Error(err),
		)
		return
	}

	// Parse timestamp
	timestamp := time.Now()
	if ts, err := strconv.ParseInt(msg.Timestamp, 10, 64); err == nil {
		timestamp = time.Unix(ts, 0)
	}

	// Create location
	locationReq := &dto.CreateLocationRequest{
		Latitude:  msg.Location.Latitude,
		Longitude: msg.Location.Longitude,
		Timestamp: &timestamp,
	}

	_, err = h.locationService.CreateLocation(
		c.Request.Context(),
		participant.ID,
		participant.EntityID,
		locationReq,
	)
	if err != nil {
		h.logger.Error("Failed to save location",
			zap.String("phone", phoneNumber),
			zap.Error(err),
		)
		return
	}

	h.logger.Info("Location saved successfully",
		zap.String("phone", phoneNumber),
		zap.String("participant_id", participant.ID.String()),
	)
}

// handleInteractiveMessage processes interactive button replies (confirmation)
func (h *WebhookHandler) handleInteractiveMessage(c *gin.Context, msg whatsapp.Message) {
	if msg.Interactive == nil || msg.Interactive.ButtonReply == nil {
		return
	}

	phoneNumber := msg.From
	buttonPayload := msg.Interactive.ButtonReply.Payload

	h.logger.Info("Received interactive reply",
		zap.String("phone", phoneNumber),
		zap.String("payload", buttonPayload),
	)

	h.processConfirmationResponse(c, phoneNumber, buttonPayload)
}

// handleButtonMessage processes button replies
func (h *WebhookHandler) handleButtonMessage(c *gin.Context, msg whatsapp.Message) {
	if msg.Button == nil {
		return
	}

	phoneNumber := msg.From
	buttonPayload := msg.Button.Payload

	h.logger.Info("Received button reply",
		zap.String("phone", phoneNumber),
		zap.String("payload", buttonPayload),
	)

	h.processConfirmationResponse(c, phoneNumber, buttonPayload)
}

// handleTextMessage processes text messages (fallback confirmation)
func (h *WebhookHandler) handleTextMessage(c *gin.Context, msg whatsapp.Message) {
	if msg.Text == nil {
		return
	}

	phoneNumber := msg.From
	text := msg.Text.Body

	h.logger.Info("Received text message",
		zap.String("phone", phoneNumber),
		zap.String("text", text),
	)

	// Simple text-based confirmation (yes/no/sim/não)
	switch text {
	case "1", "yes", "sim", "confirmo", "vou":
		h.processConfirmationResponse(c, phoneNumber, "confirm_yes")
	case "2", "no", "não", "nao", "não vou":
		h.processConfirmationResponse(c, phoneNumber, "confirm_no")
	}
}

// processConfirmationResponse processes confirmation responses
func (h *WebhookHandler) processConfirmationResponse(c *gin.Context, phoneNumber, payload string) {
	// Find participant by phone number
	participant, err := h.participantService.GetByPhoneNumber(c.Request.Context(), phoneNumber)
	if err != nil {
		h.logger.Warn("Participant not found for confirmation",
			zap.String("phone", phoneNumber),
			zap.Error(err),
		)
		return
	}

	var newStatus domain.ParticipantStatus
	switch payload {
	case "confirm_yes", "CONFIRM_YES", "yes", "1":
		newStatus = domain.ParticipantStatusConfirmed
	case "confirm_no", "CONFIRM_NO", "no", "2":
		newStatus = domain.ParticipantStatusDenied
	default:
		h.logger.Warn("Unknown confirmation payload",
			zap.String("phone", phoneNumber),
			zap.String("payload", payload),
		)
		return
	}

	// Update participant status
	err = h.participantService.UpdateStatus(c.Request.Context(), participant.ID, participant.EntityID, newStatus)
	if err != nil {
		h.logger.Error("Failed to update participant status",
			zap.String("phone", phoneNumber),
			zap.Error(err),
		)
		return
	}

	h.logger.Info("Participant confirmation processed",
		zap.String("phone", phoneNumber),
		zap.String("participant_id", participant.ID.String()),
		zap.String("status", string(newStatus)),
	)
}

// verifySignature verifies the webhook signature
func (h *WebhookHandler) verifySignature(body []byte, signature string) bool {
	if signature == "" {
		return false
	}

	// Remove "sha256=" prefix
	if len(signature) > 7 && signature[:7] == "sha256=" {
		signature = signature[7:]
	}

	mac := hmac.New(sha256.New, []byte(h.cfg.WebhookSecret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
