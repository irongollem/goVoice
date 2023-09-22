package main

import (
	"goVoice/internal/api"
	"goVoice/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	webClientAPI := api.NewWebClientAPI(cfg)
	voiceAPI := api.NewVoiceAPI(cfg)

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
