package telnyx

import (
	"context"
	"goVoice/internal/models"
	"log"

	"github.com/gin-gonic/gin"
)

/* Note: any telnyx procedure expects a 200 OK response or the event
 * will be resent over and over again.
 * The actual response to handle the event should be a separate call
 * to the telnyx API, one for every command we want to execute but always
 * at least one command should come from ...telnyx/commands.go
 */

func (t *Telnyx) answerProcedure(c *gin.Context, event Event) {
	t.answerCall(event)
}

func (t *Telnyx) startCallProcedure(c *gin.Context, event Event) {
	state, err := decodeClientState(event.Data.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
	}

	t.ConvCtrl.StartConversation(state.RulesetID, event.Data.Payload.CallControlID)
}

func (t *Telnyx) transcriptionProcedure(c *gin.Context, event Event) {
	ctx := context.Background()
	transcriptionData := event.Data.Payload.TranscriptionData
	callID := event.Data.Payload.CallControlID
	state, err := decodeClientState(event.Data.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
	}

	t.ConvCtrl.ProcessTranscription(ctx, callID, transcriptionData.Transcript, state)
}

func (t *Telnyx) hangupProcedure(c *gin.Context, event Event) {
	callID := event.Data.Payload.CallControlID
	state, err := decodeClientState(event.Data.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
	}

	ctx := context.Background()
	err = t.ConvCtrl.DB.SetConversationDone(ctx, state.RulesetID, callID)
	if err != nil {
		log.Printf("Error setting conversation done: %v", err)
		return // TODO: in its current state the conversation will be stuck in the DB
	}
	t.ConvCtrl.EndConversation(ctx, state.RulesetID, callID, state.RecordingCount)
}

func (t *Telnyx) speakStartedProcedure(c *gin.Context, event Event) {
	t.stopTranscription(event)
	t.stopRecording(event)
	log.Print("Speak started")
}

func (t *Telnyx) playbackStartedProcedure(c *gin.Context, event Event) {
	t.stopTranscription(event)
	t.stopRecording(event)
	log.Print("playback started")
}

func (t *Telnyx) speakEndedProcedure(c *gin.Context, event Event) {
	t.startTranscription(event)
	t.startRecording(event)
	log.Print("Speak ended")
}

func (t *Telnyx) playbackEndedProcedure(c *gin.Context, event Event) {
	t.startTranscription(event)
	t.startRecording(event)
	log.Print("playback ended")
}

func (t *Telnyx) recordingSavedProcedure(c *gin.Context, event Event) {
	log.Printf("Recording schema: %+v", event.Data)

	state, err := decodeClientState(event.Data.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
		return
	}
	rulesetID := state.RulesetID
	callID := event.Data.Payload.CallControlID
	url := event.Data.Payload.RecordingUrls.Mp3

	t.ConvCtrl.ProcessRecording(context.Background(), rulesetID, callID, &models.Recording{
		Url:     url,
		Purpose: state.RecordingPurpose,
	})
	if state.RecordingPurpose != state.Purpose {
		/* We update it now because recording processing can overlap. This way we make sure
		 * that the recording purpose is always set to the correct value.
		 */
		state.RecordingPurpose = state.Purpose
		t.UpdateState(callID, state)
	} else {
		log.Println("Recording came in before transcription, recording files potentially incorrectly named")
	}
}

func (t *Telnyx) recordingErrorProcedure(c *gin.Context, event Event) {
	log.Printf("Recording error: %v", event.Data.Payload.Reason)
}
