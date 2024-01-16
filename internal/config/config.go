package config

import (
	"context"
	"fmt"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/joho/godotenv"
)

type Config struct {
	ApiPort              string
	VoiceAPIAddr         string
	TelnyxAPIKey         string
	TelnyxAPIUrl         string
	TelnyxAppId          string
	GCPProjectID         string
	EmailPassword        string
	OpenAIKey            string
	OpenAIEndpoint       string
	OpenAIDeploymentName string
}

func LoadConfig() (*Config, error) {
	if os.Getenv("ENV") == "production" {
		return loadConfigFromSecretManager()
	} else {
		return LoadConfigFromEnv()
	}
}


func LoadConfigFromEnv() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
		return nil, err
	}

	return &Config{
		ApiPort:              ":" + os.Getenv("PORT"),
		TelnyxAPIKey:         os.Getenv("TELNYX_API_KEY"),
		TelnyxAPIUrl:         os.Getenv("TELNYX_API_URL"),
		TelnyxAppId:          os.Getenv("TELNYX_APP_ID"),
		GCPProjectID:         os.Getenv("GCP_PROJECT_ID"),
		EmailPassword:        os.Getenv("EMAIL_PASSWORD"),
		OpenAIKey:            os.Getenv("OPENAI_KEY"),
		OpenAIEndpoint:       os.Getenv("OPENAI_ENDPOINT"),
		OpenAIDeploymentName: os.Getenv("OPENAI_DEPLOYMENT_NAME"),
	}, nil
}

func loadConfigFromSecretManager() (*Config, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create secret manager client: %v", err)
	}
	defer client.Close()

	config := &Config{
		ApiPort:      ":" + os.Getenv("PORT"),
		GCPProjectID: os.Getenv("GCP_PROJECT_ID"),
	}

	secretNames := []string{
		"TELNYX_API_KEY",
		"TELNYX_API_URL",
		"TELNYX_APP_ID",
		"EMAIL_PASSWORD",
		"OPENAI_KEY",
		"OPENAI_ENDPOINT",
		"OPENAI_DEPLOYMENT_NAME",
	}

	for _, secretName := range secretNames {
		req := &secretmanagerpb.AccessSecretVersionRequest{
			Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", os.Getenv("GCP_PROJECT_ID"), secretName),
		}

		result, err := client.AccessSecretVersion(ctx, req)
		if err != nil {
			log.Fatalf("Failed to access secret %s: %v", secretName, err)
			return nil, err
		}

		switch secretName {
		case "TELNYX_API_KEY":
			config.TelnyxAPIKey = string(result.Payload.Data)
		case "TELNYX_API_URL":
			config.TelnyxAPIUrl = string(result.Payload.Data)
		case "TELNYX_APP_ID":
			config.TelnyxAppId = string(result.Payload.Data)
		case "EMAIL_PASSWORD":
			config.EmailPassword = string(result.Payload.Data)
		case "OPENAI_KEY":
			config.OpenAIKey = string(result.Payload.Data)
		case "OPENAI_ENDPOINT":
			config.OpenAIEndpoint = string(result.Payload.Data)
		case "OPENAI_DEPLOYMENT_NAME":
			config.OpenAIDeploymentName = string(result.Payload.Data)
		}
	}

	return config, nil
}
