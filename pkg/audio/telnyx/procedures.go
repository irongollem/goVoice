package telnyx

import (
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

func (t *Telnyx) onTranscriptionProcedure(c *gin.Context, event Event) {
	// respond to the incoming hook immediately
	c.Status(http.StatusOK)

	transcriptionData := event.Payload.TranscriptionData
	callId := event.Payload.CallControlID
	clientState := event.Payload.ClientState
	
	// await response from conversation engine
	t.ConvCtrl.ProcessTranscription(callId, transcriptionData.Transcript, clientState)
	// send response to telnyx text to speech

}