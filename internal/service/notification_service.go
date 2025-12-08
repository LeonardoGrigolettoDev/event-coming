package service

import (
	"context"
	"fmt"

	"event-coming/internal/domain"
	"event-coming/internal/whatsapp"

	"go.uber.org/zap"
)

// NotificationService define os métodos de notificação
type NotificationService interface {
	// Enviar pedido de confirmação
	SendConfirmationRequest(ctx context.Context, event *domain.Event, participant *domain.Participant) error

	// Enviar lembrete
	SendReminder(ctx context.Context, event *domain.Event, participant *domain.Participant) error

	// Enviar pedido de localização
	SendLocationRequest(ctx context.Context, event *domain.Event, participant *domain.Participant) error

	// Enviar atualização de ETA
	SendETAUpdate(ctx context.Context, event *domain.Event, participant *domain.Participant, etaMinutes int) error

	// Enviar notificação genérica
	SendMessage(ctx context.Context, phoneNumber string, message string) error
}

type notificationServiceImpl struct {
	whatsappClient *whatsapp.Client
	logger         *zap.Logger
}

func NewNotificationService(
	whatsappClient *whatsapp.Client,
	logger *zap.Logger,
) NotificationService {
	return &notificationServiceImpl{
		whatsappClient: whatsappClient,
		logger:         logger,
	}
}

// SendConfirmationRequest envia pedido de confirmação via WhatsApp
func (s *notificationServiceImpl) SendConfirmationRequest(ctx context.Context, event *domain.Event, participant *domain.Participant) error {
	message := fmt.Sprintf(
		"🎫 *Confirmação de Presença*\n\n"+
			"Olá %s!\n\n"+
			"Você está convidado para o evento:\n"+
			"📌 *%s*\n"+
			"📅 %s\n\n"+
			"Por favor, confirme sua presença respondendo:\n"+
			"✅ *SIM* - para confirmar\n"+
			"❌ *NÃO* - para recusar",
		participant.Name,
		event.Name,
		event.StartTime.Format("02/01/2006 às 15:04"),
	)

	return s.SendMessage(ctx, participant.PhoneNumber, message)
}

// SendReminder envia lembrete do evento
func (s *notificationServiceImpl) SendReminder(ctx context.Context, event *domain.Event, participant *domain.Participant) error {
	message := fmt.Sprintf(
		"⏰ *Lembrete de Evento*\n\n"+
			"Olá %s!\n\n"+
			"Seu evento está chegando:\n"+
			"📌 *%s*\n"+
			"📅 %s\n"+
			"📍 %s\n\n"+
			"Não se esqueça! 🎉",
		participant.Name,
		event.Name,
		event.StartTime.Format("02/01/2006 às 15:04"),
		getLocationAddress(event),
	)

	return s.SendMessage(ctx, participant.PhoneNumber, message)
}

// SendLocationRequest solicita a localização do participante
func (s *notificationServiceImpl) SendLocationRequest(ctx context.Context, event *domain.Event, participant *domain.Participant) error {
	message := fmt.Sprintf(
		"📍 *Compartilhe sua Localização*\n\n"+
			"Olá %s!\n\n"+
			"O evento *%s* está prestes a começar.\n\n"+
			"Por favor, compartilhe sua localização atual para calcularmos seu tempo de chegada.",
		participant.Name,
		event.Name,
	)

	return s.SendMessage(ctx, participant.PhoneNumber, message)
}

// SendETAUpdate envia atualização do tempo estimado de chegada
func (s *notificationServiceImpl) SendETAUpdate(ctx context.Context, event *domain.Event, participant *domain.Participant, etaMinutes int) error {
	var etaText string
	if etaMinutes <= 5 {
		etaText = "menos de 5 minutos"
	} else if etaMinutes <= 60 {
		etaText = fmt.Sprintf("aproximadamente %d minutos", etaMinutes)
	} else {
		hours := etaMinutes / 60
		mins := etaMinutes % 60
		etaText = fmt.Sprintf("aproximadamente %dh%02dmin", hours, mins)
	}

	// Aqui você pode enviar para o organizador do evento
	s.logger.Info("ETA Update",
		zap.String("participant", participant.Name),
		zap.Int("eta_minutes", etaMinutes),
		zap.String("eta_text", etaText),
	)

	return nil
}

// SendMessage envia mensagem genérica via WhatsApp
func (s *notificationServiceImpl) SendMessage(ctx context.Context, phoneNumber string, message string) error {
	if s.whatsappClient == nil {
		s.logger.Warn("WhatsApp client not configured, skipping message",
			zap.String("phone", phoneNumber),
		)
		return nil
	}

	s.logger.Info("Sending WhatsApp message",
		zap.String("phone", phoneNumber),
	)

	return s.whatsappClient.SendTextMessage(ctx, phoneNumber, message)
}

// getLocationAddress retorna o endereço do evento ou coordenadas
func getLocationAddress(event *domain.Event) string {
	if event.LocationAddress != nil && *event.LocationAddress != "" {
		return *event.LocationAddress
	}
	return fmt.Sprintf("%.6f, %.6f", event.LocationLat, event.LocationLng)
}
