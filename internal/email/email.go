package email

import (
	"context"
	"encoding/base64"
	"fmt"
	"goVoice/internal/config"
	"goVoice/pkg/storage"
	"log"
	"net/smtp"
	"strings"
)

type EmailProvider struct {
	password string
	storage  storage.StorageProvider
}

func NewEmailProvider(cfg *config.Config, storage storage.StorageProvider) *EmailProvider {
	return &EmailProvider{
		password: cfg.EmailPassword,
		storage:  storage,
	}
}

func (p *EmailProvider) SendEmailWithAttachment(ctx context.Context, to, subject, body string, attachments [][]byte, attachmentNames []string) error {
	from := "noreply@smartaisolutions.nl"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Mime headers
	header := make(map[string]string)
	header["From"] = from
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = `multipart/mixed; boundary="MULTIPART-MIXED-BOUNDARY"`
	header["Content-Transfer-Encoding"] = "7bit"

	// Setup message
	var message strings.Builder
	for k, v := range header {
		message.WriteString(k + ": " + v + "\n")
	}

	// Setup text
	message.WriteString("\n--MULTIPART-MIXED-BOUNDARY\n")
	message.WriteString("Content-Type: text/plain; charset=\"utf-8\"\n")
	message.WriteString("Content-Transfer-Encoding: 7bit\n")
	message.WriteString("\n" + body + "\n")

	for i, attachment := range attachments {
		// Encode attachment
		encodedRecording := base64.StdEncoding.EncodeToString(attachment)

		// Add attachment
		message.WriteString("\n--MULTIPART-MIXED-BOUNDARY\n")
		message.WriteString("Content-Type: audio/mpeg\n") // If you need to send other types of files, make this dynamic
		message.WriteString("Content-Transfer-Encoding: base64\n")
		message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\n", attachmentNames[i]))
		message.WriteString(fmt.Sprintf("\n%s\n", encodedRecording))
	}
	message.WriteString("\n--MULTIPART-MIXED-BOUNDARY--")

	auth := smtp.PlainAuth("", from, p.password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message.String()))
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	return nil
}
