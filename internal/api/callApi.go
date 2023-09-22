package api

import (
	"goVoice/internal/config"
	"net/http"

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
		voiceRoutes.GET("/", api.HandleRoot)
	}
}

func (api *VoiceAPI) HandleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from VoiceAPI",
	})
}
