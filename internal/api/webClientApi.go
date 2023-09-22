package api

import (
	"goVoice/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebClientAPI struct {
	Router *gin.Engine
}

func NewWebClientAPI(cfg *config.Config) *WebClientAPI {
	ginEngine:= gin.Default()
	api := &WebClientAPI{ Router: ginEngine }
	api.routes()
	return api
}

func (api *WebClientAPI) routes() {
		api.Router.GET("/", api.HandleRoot)
}

func (api *WebClientAPI) HandleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from WebClientAPI",
	})
}
