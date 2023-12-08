package telnyx

import (
	"context"
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
	log.Printf("Answering call: %v", event)

	t.answerCall(event)
}

func (t *Telnyx) startCallProcedure(c *gin.Context, event Event) {
	// respond to the incoming hook immediately
	c.Status(http.StatusOK)
	log.Printf("Call started: %v", event)

	// Start recording the call and start transcription
	// Any further process should only start once we know transcription
	// has been started
	t.startTranscription(event)
	t.startRecording(event)

	t.ConvCtrl.StartConversation(event.Data.Payload.CallControlID)
}

func (t *Telnyx) transcriptionProcedure(c *gin.Context, event Event) {
	// respond to the incoming hook immediately
	c.Status(http.StatusOK)
	log.Printf("Transcription received: %v", event)

	ctx := context.Background()
	transcriptionData := event.Data.Payload.TranscriptionData
	callId := event.Data.Payload.CallControlID
	state, err := decodeClientState(event.Data.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
	}

	t.ConvCtrl.ProcessTranscription(ctx, callId, transcriptionData.Transcript, state)
}

func (t *Telnyx) hangupProcedure(c *gin.Context, event Event) {
	// respond to the incoming hook immediately
	c.Status(http.StatusOK)
	log.Printf("Call ended: %v", event)

	callId := event.Data.Payload.CallControlID
	state, err := decodeClientState(event.Data.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
	}

	ctx := c.Request.Context()
	// TODO check if recording is done, if not, wait for it to finish
	t.ConvCtrl.EndConversation(ctx, state.RulesetID, callId)
}

/**
* For now we don't do anything with the speak started event but are
* required to respond to the incoming hook
*/
func (t *Telnyx) speakStartedProcedure(c *gin.Context, event Event) {
	// respond to the incoming hook immediately
	c.Status(http.StatusOK)
	log.Print("Speak started")
}

/**
* For now we don't do anything with the speak ended event but are
* required to respond to the incoming hook
*/
func (t *Telnyx) speakEndedProcedure(c *gin.Context, event Event) {
	// respond to the incoming hook immediately
	c.Status(http.StatusOK)
	log.Print("Speak ended")
}

func (t *Telnyx) recordingSavedProcedure(c *gin.Context, event Event) {
	// respond to the incoming hook immediately
	c.Status(http.StatusOK)
	log.Printf("Recording saved: %v", event.Data.RecordingId)
	
	callId := event.Data.Payload.CallControlID
	recordingId := event.Data.RecordingId
	state, err := decodeClientState(event.Data.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
	}

	ctx := c.Request.Context()
	recChan, errChan := t.GetRecording(recordingId)
	select {
		case telnyxRecording := <- recChan:
			genericRecording := convertToRecording(telnyxRecording)
			t.ConvCtrl.ProcessRecording(ctx, state.RulesetID, callId, genericRecording)
		case err := <- errChan:
			log.Printf("Error getting recordings: %v", err)
		}
}

func (t *Telnyx) recordingErrorProcedure(c *gin.Context, event Event) {
	// respond to the incoming hook immediately
	c.Status(http.StatusOK)
	log.Printf("Recording error: %v", event.Data.Payload.Reason)
}