package api

import (
	"goVoice/internal/app/conversation"
	"goVoice/internal/config"
	"goVoice/internal/email"
	"goVoice/pkg/ai"
	"goVoice/pkg/audio"
	"goVoice/pkg/audio/telnyx"
	"goVoice/pkg/db"
	"goVoice/pkg/storage"

	"github.com/gin-gonic/gin"
)

type VoiceAPI struct {
	Router *gin.Engine
}

func NewVoiceAPI(cfg *config.Config, storage storage.StorageProvider, db db.DbProvider, ai ai.AIProvider, router *gin.Engine) *VoiceAPI {
	api := &VoiceAPI{Router: router}

	// We can replace the Telnyx struct with any other provider
	// as long as they implement the CallProvider interface
	convCtrl := &conversation.Controller{
		Storage: storage,
		DB:      db,
		AI:      ai,
		Email:   email.NewEmailProvider(cfg),
	}
	client := telnyx.NewTelnyxClient(cfg, convCtrl)
	// client.SetBucketCredentials(cfg)
	convCtrl.CallProvider = client

	api.routes(client)
	return api
}

func (api *VoiceAPI) routes(client audio.CallProvider) {
	voiceRoutes := api.Router.Group("/call")
	{
		voiceRoutes.GET("/", client.IAmLive)
		voiceRoutes.POST("/", client.HandleWebHook)
	}
}
