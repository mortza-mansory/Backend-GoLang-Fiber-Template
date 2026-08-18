// Package email provides a generic outbound email boundary. Real provider
// integrations (SMTP, SES, SendGrid, etc.) should implement the Sender
// interface; this template ships with an SMTP client and an optional no-op.
package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/yourorg/go-fiber-template/internal/config"
)

// Message is a minimal outbound email.
type Message struct {
	To      []string
	Subject string
	Body    string
}

// Sender sends email messages.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// SMTPClient sends emails via an SMTP server.
type SMTPClient struct {
	cfg  config.EmailService
	auth smtp.Auth
}

// NewSMTPClient creates an SMTPClient from config.
func NewSMTPClient(cfg config.EmailService) *SMTPClient {
	var auth smtp.Auth
	if cfg.User != "" {
		auth = smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	}
	return &SMTPClient{cfg: cfg, auth: auth}
}

// Send delivers msg via SMTP.
func (s *SMTPClient) Send(ctx context.Context, msg Message) error {
	if len(msg.To) == 0 {
		return fmt.Errorf("email: no recipients")
	}
	body := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s",
		msg.To[0], msg.Subject, msg.Body)
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	return smtp.SendMail(addr, s.auth, s.cfg.From, msg.To, []byte(body))
}

// NoopSender is a placeholder that logs and discards. Useful when SMTP is not
// configured (e.g. local development) without failing the app.
type NoopSender struct{}

// NewNoopSender creates a NoopSender.
func NewNoopSender() *NoopSender { return &NoopSender{} }

// Send does nothing. It exists so wiring can always provide a Sender.
func (NoopSender) Send(ctx context.Context, msg Message) error {
	return nil
}
