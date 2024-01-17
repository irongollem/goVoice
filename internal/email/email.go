package email

import (
	"context"
	"goVoice/internal/config"
	"io"
	"log"

	"gopkg.in/gomail.v2"
)

type EmailProvider struct {
	password string
}

func NewEmailProvider(cfg *config.Config) *EmailProvider {
	return &EmailProvider{
		password: cfg.EmailPassword,
	}
}

func (p *EmailProvider) SendEmailWithAttachment(ctx context.Context, to, subject, body string, attachments [][]byte, attachmentNames []string) error {
	from := "noreply@smartaisolutions.nl"
	smtpHost := "smtp.gmail.com"

	smtpPort := 587

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)


	for i, attachment := range attachments {
		m.Attach((attachmentNames[i]), gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(attachment)
			return err
		}))
	}

	d := gomail.NewDialer(smtpHost, smtpPort, from, p.password)
	
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	return nil
}
