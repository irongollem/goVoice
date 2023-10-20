package audio

import "github.com/gin-gonic/gin"

type CallProvider interface {
	HandleWebHook(c *gin.Context)
	SpeakText(callId string, text string, clientState *string) (chan bool, chan error)
	EndCall(callId string) (chan bool, chan error)
}