package internal

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	config *Config
}

func NewEmailService(config *Config) *EmailService {
	return &EmailService{
		config: config,
	}
}

// SendPinEmail sends a PIN to the user's email
func (e *EmailService) SendPinEmail(to, pin string, expiryMins int) error {
	subject := "Your Authentication PIN"
	body := fmt.Sprintf(`
<html>
<body>
	<h2>Authentication PIN</h2>
	<p>Your authentication PIN is:</p>
	<h1 style="letter-spacing: 5px; font-size: 36px;">%s</h1>
	<p>This PIN will expire in %d minutes.</p>
	<p>If you didn't request this PIN, please ignore this email.</p>
</body>
</html>
	`, pin, expiryMins)

	return e.sendEmail(to, subject, body)
}

// sendEmail sends an email using SMTP
func (e *EmailService) sendEmail(to, subject, body string) error {
	from := e.config.EmailFrom
	password := e.config.SMTPPassword

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Setup authentication
	auth := smtp.PlainAuth("", e.config.SMTPUsername, password, e.config.SMTPHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", e.config.SMTPHost, e.config.SMTPPort)
	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
