package telnyx

import (
	"goVoice/internal/app/conversation"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Telnyx struct {
	APIKey      string
	CommandPath string
	ConvCtrl		*conversation.Controller
}

func NewTelnyxClient(apiKey string, apiUrl string, convCtrl *conversation.Controller) *Telnyx {
	return &Telnyx{
		APIKey:      apiKey,
		CommandPath: apiUrl + "/calls/",
		ConvCtrl:    convCtrl,
	}
}

func (t *Telnyx) HandleWebHook(c *gin.Context) {
	c.Writer.Header().Set("Authorization", "Bearer "+t.APIKey)

	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		log.Printf("Error decoding webhook data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error decoding webhook data",
		})
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	callType := event.EventType

	switch callType {
	case "call.initiated":
		t.answerProcedure(c, event)
	case "call.answered":
		t.startCallProcedure(c, event)
	case "call.hangup":
		c.AbortWithStatus(http.StatusNotFound)
	case "call.speak.ended":
		c.AbortWithStatus(http.StatusNotFound)
	case "call.speak.started":
		c.AbortWithStatus(http.StatusNotFound)
	case "call.recording.saved":
		c.AbortWithStatus(http.StatusNotFound)
	case "call.transcription":
		t.onTranscriptionProcedure(c, event)
		// When a call is answered using a socket, its content will be streamed
	case "streaming.started":
		c.AbortWithStatus(http.StatusNotFound)
	case "streaming.stopped":
		c.AbortWithStatus(http.StatusNotFound)
	}
}
