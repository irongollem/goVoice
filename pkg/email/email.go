package email

import (
	"context"
	"goVoice/internal/config"
	"goVoice/internal/models"
	// "goVoice/pkg/email/sendgridClient"

	"io"
)

type EmailProvider interface {
	SendConversationEmail(ctx context.Context, recording io.ReadCloser, responses []models.ConversationStepResponse) error
}

func newEmailHandler(cfg *config.Config) (EmailProvider, error) {
	panic("implement me")
}
