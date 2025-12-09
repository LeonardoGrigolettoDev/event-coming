package whatsapp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"event-coming/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	cfg := &config.WhatsAppConfig{
		BaseURL:       "https://graph.facebook.com",
		APIVersion:    "v18.0",
		PhoneNumberID: "123456789",
		AccessToken:   "test-token",
	}

	client := NewClient(cfg)

	assert.NotNil(t, client)
	assert.Equal(t, "https://graph.facebook.com/v18.0/123456789", client.baseURL)
	assert.NotNil(t, client.httpClient)
}

func TestClient_SendTemplateMessage(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful send",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:        "server error",
			statusCode:  http.StatusInternalServerError,
			wantErr:     true,
			errContains: "unexpected status code: 500",
		},
		{
			name:        "unauthorized",
			statusCode:  http.StatusUnauthorized,
			wantErr:     true,
			errContains: "unexpected status code: 401",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Contains(t, r.URL.Path, "/messages")
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Contains(t, r.Header.Get("Authorization"), "Bearer")

				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			cfg := &config.WhatsAppConfig{
				BaseURL:       server.URL,
				APIVersion:    "v18.0",
				PhoneNumberID: "123456789",
				AccessToken:   "test-token",
			}

			client := NewClient(cfg)

			req := &TemplateMessageRequest{
				MessagingProduct: "whatsapp",
				RecipientType:    "individual",
				To:               "+5511999999999",
				Type:             "template",
				Template: Template{
					Name:     "test_template",
					Language: Language{Code: "en"},
				},
			}

			err := client.SendTemplateMessage(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_SendConfirmationRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var req TemplateMessageRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, "whatsapp", req.MessagingProduct)
		assert.Equal(t, "+5511999999999", req.To)
		assert.Equal(t, "event_confirmation", req.Template.Name)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.WhatsAppConfig{
		BaseURL:       server.URL,
		APIVersion:    "v18.0",
		PhoneNumberID: "123456789",
		AccessToken:   "test-token",
	}

	client := NewClient(cfg)

	err := client.SendConfirmationRequest(
		context.Background(),
		"+5511999999999",
		"John Doe",
		"Test Event",
		time.Now().Add(24*time.Hour),
	)

	assert.NoError(t, err)
}

func TestClient_SendLocationRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req TemplateMessageRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, "whatsapp", req.MessagingProduct)
		assert.Equal(t, "+5511999999999", req.To)
		assert.Equal(t, "location_request", req.Template.Name)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.WhatsAppConfig{
		BaseURL:       server.URL,
		APIVersion:    "v18.0",
		PhoneNumberID: "123456789",
		AccessToken:   "test-token",
	}

	client := NewClient(cfg)

	err := client.SendLocationRequest(
		context.Background(),
		"+5511999999999",
		"John Doe",
		"Test Event",
	)

	assert.NoError(t, err)
}

func TestClient_SendTextMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, "whatsapp", req["messaging_product"])
		assert.Equal(t, "+5511999999999", req["to"])
		assert.Equal(t, "text", req["type"])

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.WhatsAppConfig{
		BaseURL:       server.URL,
		APIVersion:    "v18.0",
		PhoneNumberID: "123456789",
		AccessToken:   "test-token",
	}

	client := NewClient(cfg)

	err := client.SendTextMessage(
		context.Background(),
		"+5511999999999",
		"Hello, this is a test message!",
	)

	assert.NoError(t, err)
}

func TestClient_SendTextMessage_NetworkError(t *testing.T) {
	cfg := &config.WhatsAppConfig{
		BaseURL:       "http://localhost:1", // Invalid port
		APIVersion:    "v18.0",
		PhoneNumberID: "123456789",
		AccessToken:   "test-token",
	}

	client := NewClient(cfg)

	err := client.SendTextMessage(
		context.Background(),
		"+5511999999999",
		"Hello",
	)

	assert.Error(t, err)
}

func TestTemplateMessageRequest_JSON(t *testing.T) {
	req := TemplateMessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               "+5511999999999",
		Type:             "template",
		Template: Template{
			Name:     "test_template",
			Language: Language{Code: "en"},
			Components: []Component{
				{
					Type: "body",
					Parameters: []Parameter{
						{Type: "text", Text: "John"},
						{Type: "text", Text: "Event Name"},
					},
				},
			},
		},
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded TemplateMessageRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, req.To, decoded.To)
	assert.Equal(t, req.Template.Name, decoded.Template.Name)
}
