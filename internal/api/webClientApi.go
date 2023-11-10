package api

import (
	"goVoice/internal/config"
	"goVoice/pkg/db"
	"goVoice/pkg/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebClientAPI struct {
	Router  *gin.Engine
	storage storage.StorageProvider
	db      db.DbProvider
}

func NewWebClientAPI(cfg *config.Config, storageHandler storage.StorageProvider, dbHandler db.DbProvider, router *gin.Engine) *WebClientAPI {
	api := &WebClientAPI{Router: router, storage: storageHandler, db: dbHandler}
	api.routes()
	return api
}

func (api *WebClientAPI) routes() {
	api.Router.GET("/", api.HandleRoot)
}

func (api *WebClientAPI) HandleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from WebClientAPI, it seems your API is running!",
	})
}
