package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WebClientAPIAddr   string
	VoiceAPIAddr       string
	TelnyxAPIKey       string
	TelnyxAPIUrl       string
	TelnyxAppId				 string
	GCPCredentialsFile string
	GCPProjectID       string
	SendgridAPIKey     string
	EmailPassword			 string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
		return nil, err
	}

	return &Config{
		WebClientAPIAddr:   ":8080",
		VoiceAPIAddr:       ":8081",
		TelnyxAPIKey:       os.Getenv("TELNYX_API_KEY"),
		TelnyxAPIUrl:       os.Getenv("TELNYX_API_URL"),
		TelnyxAppId:        os.Getenv("TELNYX_APP_ID"),
		GCPCredentialsFile: os.Getenv("GCP_CREDENTIALS_FILE"),
		GCPProjectID:       os.Getenv("GCP_PROJECT_ID"),
		SendgridAPIKey:     os.Getenv("SENDGRID_API_KEY"),
		EmailPassword:      os.Getenv("EMAIL_PASSWORD"),
	}, nil
}
