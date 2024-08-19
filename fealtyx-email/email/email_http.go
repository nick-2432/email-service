package email

// SendEmailRequest represents the expected request body for sending an email
type SendEmailRequest struct {
	Content string `json:"content"`
}

// SendEmailResponse represents the response body for sending an email
type SendEmailResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
