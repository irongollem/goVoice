package audio

import (
	"goVoice/internal/models"

	"github.com/gin-gonic/gin"
)

type CallProvider interface {
	HandleWebHook(c *gin.Context)
	IAmLive(c *gin.Context)
	SpeakText(callID string, text string, clientState *models.ClientState) (chan bool, chan error)
	PlayAudioUrl(callID string, audioUrl string, clientState *models.ClientState) (chan bool, chan error)
	GetRecordingMp3(recording *models.Recording) (chan []byte, chan error)
	EndCall(callID string) (chan bool, chan error)
}
