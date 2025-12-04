package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"event-coming/internal/config"
)

// Client handles WhatsApp Cloud API interactions
type Client struct {
	config     *config.WhatsAppConfig
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new WhatsApp client
func NewClient(cfg *config.WhatsAppConfig) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: fmt.Sprintf("%s/%s/%s", cfg.BaseURL, cfg.APIVersion, cfg.PhoneNumberID),
	}
}

// SendTemplateMessage sends a template message
func (c *Client) SendTemplateMessage(ctx context.Context, req *TemplateMessageRequest) error {
	url := fmt.Sprintf("%s/messages", c.baseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.AccessToken))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// SendConfirmationRequest sends a confirmation request to a participant
func (c *Client) SendConfirmationRequest(ctx context.Context, phoneNumber, participantName, eventName string, eventTime time.Time) error {
	req := &TemplateMessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phoneNumber,
		Type:             "template",
		Template: Template{
			Name:     "event_confirmation",
			Language: Language{Code: "en"},
			Components: []Component{
				{
					Type: "body",
					Parameters: []Parameter{
						{Type: "text", Text: participantName},
						{Type: "text", Text: eventName},
						{Type: "text", Text: eventTime.Format("2006-01-02 15:04")},
					},
				},
			},
		},
	}

	return c.SendTemplateMessage(ctx, req)
}

// SendLocationRequest sends a location request to a participant
func (c *Client) SendLocationRequest(ctx context.Context, phoneNumber, participantName, eventName string) error {
	req := &TemplateMessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phoneNumber,
		Type:             "template",
		Template: Template{
			Name:     "location_request",
			Language: Language{Code: "en"},
			Components: []Component{
				{
					Type: "body",
					Parameters: []Parameter{
						{Type: "text", Text: participantName},
						{Type: "text", Text: eventName},
					},
				},
			},
		},
	}

	return c.SendTemplateMessage(ctx, req)
}
