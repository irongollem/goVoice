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
	GCPCredentialsFile string
	SendgridAPIKey     string
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
		GCPCredentialsFile: os.Getenv("GCP_CREDENTIALS_FILE"),
		SendgridAPIKey:     os.Getenv("SENDGRID_API_KEY"),
	}, nil
}
