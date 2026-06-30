package mail

import (
	"context"
	"log/slog"
)

type MailService interface {
	SendEmail(ctx context.Context, to string, subject string, body string) error
}

type mockMailService struct{}

func NewMockMailService() MailService {
	return &mockMailService{}
}

func (m *mockMailService) SendEmail(ctx context.Context, to string, subject string, body string) error {
	slog.Info("Sending email asynchronously (Mock)", "to", to, "subject", subject, "body", body)
	return nil
}
