package whatsapp

import "time"

// WebhookPayload represents the webhook payload from WhatsApp
type WebhookPayload struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

// Entry represents a webhook entry
type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

// Change represents a webhook change
type Change struct {
	Value Value  `json:"value"`
	Field string `json:"field"`
}

// Value represents the change value
type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts,omitempty"`
	Messages         []Message `json:"messages,omitempty"`
	Statuses         []Status  `json:"statuses,omitempty"`
}

// Metadata represents message metadata
type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

// Contact represents a contact
type Contact struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

// Profile represents a contact profile
type Profile struct {
	Name string `json:"name"`
}

// Message represents a WhatsApp message
type Message struct {
	From      string         `json:"from"`
	ID        string         `json:"id"`
	Timestamp string         `json:"timestamp"`
	Type      string         `json:"type"`
	Text      *TextContent   `json:"text,omitempty"`
	Location  *Location      `json:"location,omitempty"`
	Button    *ButtonReply   `json:"button,omitempty"`
	Interactive *InteractiveReply `json:"interactive,omitempty"`
}

// TextContent represents text message content
type TextContent struct {
	Body string `json:"body"`
}

// Location represents a location message
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// ButtonReply represents a button reply
type ButtonReply struct {
	Payload string `json:"payload"`
	Text    string `json:"text"`
}

// InteractiveReply represents an interactive reply
type InteractiveReply struct {
	Type        string       `json:"type"`
	ButtonReply *ButtonReply `json:"button_reply,omitempty"`
}

// Status represents a message status update
type Status struct {
	ID           string       `json:"id"`
	Status       string       `json:"status"`
	Timestamp    string       `json:"timestamp"`
	RecipientID  string       `json:"recipient_id"`
	Conversation Conversation `json:"conversation,omitempty"`
	Pricing      Pricing      `json:"pricing,omitempty"`
}

// Conversation represents conversation info
type Conversation struct {
	ID                 string    `json:"id"`
	Origin             Origin    `json:"origin"`
	ExpirationTimestamp string   `json:"expiration_timestamp,omitempty"`
}

// Origin represents conversation origin
type Origin struct {
	Type string `json:"type"`
}

// Pricing represents pricing info
type Pricing struct {
	Billable     bool   `json:"billable"`
	PricingModel string `json:"pricing_model"`
	Category     string `json:"category"`
}

// ParseTimestamp parses WhatsApp timestamp
func ParseTimestamp(ts string) (time.Time, error) {
	return time.Parse(time.RFC3339, ts)
}
