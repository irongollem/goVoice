package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WebClientAPIAddr string
	VoiceAPIAddr     string
	plivoAuthID      string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
		return nil, err
	}




	return &Config{
		WebClientAPIAddr: ":8080",
		VoiceAPIAddr:     ":8081",
		plivoAuthID:     os.Getenv("PLIVO_AUTH_ID") ,

	}, nil
}
