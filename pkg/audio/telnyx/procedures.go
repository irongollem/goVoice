package telnyx

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

/* Note: any telnyx procedure expects a 200 OK response or the event
 * will be resent over and over again.
 * The actual response to handle the event should be a separate call
 * to the telnyx API, one for every command we want to execute but always
 * at least one command should come from ...telnyx/commands.go
 */

func (t *Telnyx) answerProcedure(c *gin.Context, event Event) {
	// respond to the incoming hook immediately
	c.Status(http.StatusOK)

	t.answerCall(event)
}

func (t *Telnyx) startCallProcedure(c *gin.Context, event Event) {
	// respond to the incoming hook immediately
	c.Status(http.StatusOK)

	// Start recording the call and start transcription
	// Any further process should only start once we know transcription
	// has been started
	t.startTranscription(event)
	t.startRecording(event)

	t.ConvCtrl.StartConversation(event.Payload.CallControlID) // TODO: check if callControlId is the correct one
}

func (t *Telnyx) transcriptionProcedure(c *gin.Context, event Event) {
	// respond to the incoming hook immediately
	c.Status(http.StatusOK)
	ctx := c.Request.Context()

	transcriptionData := event.Payload.TranscriptionData
	callId := event.Payload.CallControlID
	state, err := decodeClientState(event.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
	}

	t.ConvCtrl.ProcessTranscription(ctx, callId, transcriptionData.Transcript, state)
}

func (t *Telnyx) hangupProcedure(c *gin.Context, event Event) {
	// respond to the incoming hook immediately
	c.Status(http.StatusOK)

	callId := event.Payload.CallControlID
	state, err := decodeClientState(event.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
	}

	ctx := c.Request.Context()

	t.ConvCtrl.EndConversation(ctx, state.RulesetID, callId)
}
