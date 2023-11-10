package api

import (
	"goVoice/internal/app/conversation"
	"goVoice/internal/config"
	"goVoice/pkg/audio"
	"goVoice/pkg/audio/telnyx"
	"goVoice/pkg/db"
	"goVoice/pkg/storage"

	"github.com/gin-gonic/gin"
)

type VoiceAPI struct {
	Router *gin.Engine
}

func NewVoiceAPI(cfg *config.Config, storage storage.StorageProvider, db db.DbProvider, router *gin.Engine) *VoiceAPI {
	api := &VoiceAPI{Router: router}

	// We can replace the Telnyx struct with any other provider
	// as long as they implement the CallProvider interface
	convCtrl := &conversation.Controller{
		Storage: storage,
		DB:      db,
	}
	client := telnyx.NewTelnyxClient(cfg, convCtrl)
	client.SetBucketCredentials(cfg)
	convCtrl.Provider = client

	api.routes(client)
	return api
}

func (api *VoiceAPI) routes(client audio.CallProvider) {
	voiceRoutes := api.Router.Group("/call")
	{
		voiceRoutes.POST("/", client.HandleWebHook)
	}
}
