package helpers

import (
	"fmt"

	"github.com/instrlabs/shared/functionx"
)

// EmailSender interface for email operations
type EmailSender interface {
	SendEmail(to, subject, body string)
}

// EmailService implements EmailSender
type EmailService struct{}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	return &EmailService{}
}

// SendEmail sends an email using the shared functionx utility
func (e *EmailService) SendEmail(to, subject, body string) {
	functionx.SendEmail(to, subject, body)
}

// SendPinEmail sends a PIN email
func (e *EmailService) SendPinEmail(email, pin string) {
	subject := "Your Login PIN"
	body := fmt.Sprintf("Your one-time PIN is: %s. It expires in 10 minutes.", pin)
	e.SendEmail(email, subject, body)
}
