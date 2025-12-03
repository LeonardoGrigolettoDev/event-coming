package whatsapp

// TemplateMessageRequest represents a template message request
type TemplateMessageRequest struct {
	MessagingProduct string   `json:"messaging_product"`
	RecipientType    string   `json:"recipient_type"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}

// Template represents a message template
type Template struct {
	Name       string      `json:"name"`
	Language   Language    `json:"language"`
	Components []Component `json:"components,omitempty"`
}

// Language represents template language
type Language struct {
	Code string `json:"code"`
}

// Component represents a template component
type Component struct {
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters,omitempty"`
}

// Parameter represents a template parameter
type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// InteractiveMessage represents an interactive message
type InteractiveMessage struct {
	MessagingProduct string      `json:"messaging_product"`
	RecipientType    string      `json:"recipient_type"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Interactive      Interactive `json:"interactive"`
}

// Interactive represents interactive content
type Interactive struct {
	Type   string  `json:"type"`
	Body   Body    `json:"body"`
	Action Action  `json:"action"`
}

// Body represents message body
type Body struct {
	Text string `json:"text"`
}

// Action represents interactive action
type Action struct {
	Buttons []Button `json:"buttons,omitempty"`
}

// Button represents an interactive button
type Button struct {
	Type  string `json:"type"`
	Reply Reply  `json:"reply"`
}

// Reply represents button reply
type Reply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
