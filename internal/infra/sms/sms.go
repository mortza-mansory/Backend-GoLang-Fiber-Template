// Package sms provides a generic outbound SMS boundary. Real provider
// integrations (Twilio, etc.) should implement the Sender interface. This
// template ships with a no-op placeholder.
package sms

import "context"

// Message is a minimal outbound SMS.
type Message struct {
	To   string
	Body string
}

// Sender sends SMS messages.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// NoopSender is a placeholder that discards messages.
type NoopSender struct{}

// NewNoopSender creates a NoopSender.
func NewNoopSender() *NoopSender { return &NoopSender{} }

// Send does nothing.
func (NoopSender) Send(ctx context.Context, msg Message) error {
	return nil
}
