package telnyx

import (
	"bytes"
	"goVoice/internal/app/conversation"
	"goVoice/internal/config"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

type Telnyx struct {
	APIKey   string
	APIurl   url.URL
	ConvCtrl *conversation.Controller
}

func NewTelnyxClient(cfg *config.Config, convCtrl *conversation.Controller) *Telnyx {
	apiUrl, err := url.Parse(cfg.TelnyxAPIUrl)
	if err != nil {
		log.Fatalf("Error parsing Telnyx API URL: %v", err)
	}
	
	client := &Telnyx{
		APIKey:   cfg.TelnyxAPIKey,
		APIurl:   *apiUrl,
		ConvCtrl: convCtrl,
	}
	client.setBucketCredentials(cfg)
	return client
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
	// case "streaming.started":
	// 	// When a call is answered using a socket, its content will be streamed
	// 	c.AbortWithStatus(http.StatusNotFound)
	// case "streaming.stopped":
	// 	c.AbortWithStatus(http.StatusNotFound)
	default:
		log.Printf("Unknown event type: %s received", callType)
		c.Status(http.StatusOK)
	}
}

func (t *Telnyx) setBucketCredentials(cfg *config.Config) error {
	url := `https://api.telnyx.com/v2/custom_storage_credentials/` + cfg.TelnyxAppId
	credentials, err := os.ReadFile(cfg.GCPCredentialsFile)
	if err != nil {
		log.Printf("Error reading GCP credentials file: %v", err)
		return err
	}
	payload := `{
		"backend": "gcs",
		"configuration": {
			"bucket": "govoice-recordings",
			"credentials": ` + string(credentials) + `,
		}
	}`

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
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
