package api

import (
	"encoding/json"
	"goVoice/internal/config"
	"goVoice/internal/models"
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
	api.Router.POST("/ruleset", api.apiKeyRequired(), api.HandleRulesetUpload)
}

func (api *WebClientAPI) apiKeyRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		const apiKeyHeader = "X-API-KEY"
		const apiKey = "58E1A0D885E1CDB3F48477DD4140FC2CBA10D4D2"
		// TODO store api keys in database, delete this one as its in the codebase history for ever
		// and give every client their own api key

		if c.GetHeader(apiKeyHeader) != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (api *WebClientAPI) HandleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from WebClientAPI, it seems your API is running!",
	})
}

func (api *WebClientAPI) HandleRulesetUpload(c *gin.Context) {
	// This function takes the JSON body, marshals it to models.ConversationRuleSet. and stores it in the database using the dbHandler if it is correctly formatted.
	
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error reading request body",
		})
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var ruleset models.ConversationRuleSet
	err = json.Unmarshal(jsonData, &ruleset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error parsing request body as ruleset",
		})
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	context := c.Request.Context()
	err = api.db.AddRuleset(context, &ruleset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error writing ruleset to database",
		})
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Ruleset successfully uploaded",
	})
}
