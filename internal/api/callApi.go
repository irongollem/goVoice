package api

import (
	"goVoice/internal/config"
	"goVoice/pkg/audio/plivo"

	"github.com/gin-gonic/gin"
)

type VoiceAPI struct {
	Router *gin.Engine
}

func NewVoiceAPI(cfg *config.Config) *VoiceAPI {
	ginEngine := gin.Default()
	api := &VoiceAPI{Router: ginEngine}
	api.routes()
	return api
}

func (api *VoiceAPI) routes() {
	voiceRoutes := api.Router.Group("/voice")
	{
		voiceRoutes.GET("/start_stream", plivo.StartStream)
		voiceRoutes.GET("/ws", plivo.StartStreamSocket)
	}
}

