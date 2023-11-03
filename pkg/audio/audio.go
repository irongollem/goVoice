package audio

import (
	"goVoice/internal/models"

	"github.com/gin-gonic/gin"
)

type CallProvider interface {
	HandleWebHook(c *gin.Context)
	SpeakText(callId string, text string, clientState *models.ClientState) (chan bool, chan error)
	EndCall(callId string) (chan bool, chan error)
}