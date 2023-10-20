package main

import (
	"goVoice/internal/api"
	"goVoice/internal/config"
	"goVoice/pkg/storage"
	"goVoice/pkg/db"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	storageHandler, err := storage.NewStorageHandler(cfg)
	dbHandler, err := db.NewDbHandler(cfg)
	//TODO handle the errors

	// Create the API for the UI
	webClientAPI := api.NewWebClientAPI(cfg, &storageHandler, dbHandler)
	// Create the API for the call manager
	voiceAPI := api.NewVoiceAPI(cfg, &storageHandler, dbHandler)
	
	if err != nil {
		log.Fatalf("Failed to create storage handler: %v", err)
	}

	go func() {
		if err := webClientAPI.Router.Run(cfg.WebClientAPIAddr); err != nil {
			log.Fatalf("Failed to start web client server: %v", err)
		}
	}()

	go func() {
		if err := voiceAPI.Router.Run(cfg.VoiceAPIAddr); err != nil {
			log.Fatalf("Failed to start call API server: %v", err)
		}
	}()
}
