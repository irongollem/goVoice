package main

import (
	"fmt"
	"goVoice/internal/api"
	"goVoice/internal/config"
	"goVoice/pkg/db"
	"goVoice/pkg/host"
	"goVoice/pkg/storage"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println(os.Getwd())
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	storageHandler, err := storage.NewStorageHandler(cfg)
	if err != nil {
		log.Fatalf("Failed to create storage handler: %v", err)
	}
	dbHandler, err := db.InitiateDBClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create db handler: %v", err)
	}

	router := gin.Default()
	router.Use(host.ValidateXForwardedFor)
	// Create the API for the UI
	api.NewWebClientAPI(cfg, storageHandler, dbHandler, router)
	// Create the API for the call manager
	api.NewVoiceAPI(cfg, storageHandler, dbHandler, router)
	if err := router.Run(cfg.ApiPort); err != nil {
		log.Fatalf("Failed to start web client server: %v", err)
	}
}
