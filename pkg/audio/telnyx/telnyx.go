package telnyx

import (
	"bytes"
	"encoding/json"
	"goVoice/internal/app/conversation"
	"goVoice/internal/config"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Telnyx struct {
	APIKey   string
	APIUrl   *url.URL
	ConvCtrl *conversation.Controller
}

func NewTelnyxClient(cfg *config.Config, convCtrl *conversation.Controller) *Telnyx {
	apiUrl, err := url.Parse(cfg.TelnyxAPIUrl)
	if err != nil {
		log.Fatalf("Error parsing Telnyx API URL: %v", err)
	}

	client := &Telnyx{
		APIKey:   cfg.TelnyxAPIKey,
		APIUrl:   apiUrl,
		ConvCtrl: convCtrl,
	}
	return client
}

func (t *Telnyx) IAmLive (c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from Telnyx, it seems your API is running!",
	})
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

	callType := event.Data.EventType

	// respond to the incoming hook immediately, then process the event in the background
	c.Status(http.StatusOK)
	log.Printf("Event type: %s received; starting procedure", callType)
	go func() {
		switch callType {
		case "call.initiated":
			t.answerProcedure(c, event)
		case "call.answered":
			t.startCallProcedure(c, event)
		case "call.hangup":
			t.hangupProcedure(c, event)
		case "call.speak.ended":
			t.speakEndedProcedure(c, event)
		case "call.speak.started":
			t.speakStartedProcedure(c, event)
		case "call.recording.saved":
			t.recordingSavedProcedure(c, event)
		case "call.recording.error":
			t.recordingErrorProcedure(c, event)
		case "call.transcription":
			t.transcriptionProcedure(c, event)
		case "call.playback.started":
			t.playbackStartedProcedure(c, event)
		case "call.playback.ended":
			t.playbackEndedProcedure(c, event)
		default:
			log.Printf("Unknown event type: %s received", callType)
		}
	}()
}

func (t *Telnyx) SetBucketCredentials(cfg *config.Config) error {
	url := `https://api.telnyx.com/v2/custom_storage_credentials/` + cfg.TelnyxAppId

	payload := CredentialsPayload{
		Backend: "gcs",
		Configuration: CredentialsConfiguration{
			Bucket:      "govoice-recordings",
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return err
	}

	// TODO: temporarily turned off custom buckets
	// t.sendCommand("POST", &payload, "custom_storage_credentials", cfg.TelnyxAppId)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.APIKey)

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending command to store on GCP", err)
		return err
	}

	return nil
}
