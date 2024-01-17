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
	t.startRecording(event)
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
	log.Printf("transcription received for %s", state.Purpose)

	t.ConvCtrl.ProcessTranscription(ctx, callID, transcriptionData.Transcript, state)
}

func (t *Telnyx) hangupProcedure(c *gin.Context, event Event) {
	callID := event.Data.Payload.CallControlID
	state, err := decodeClientState(event.Data.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
	}
	t.stopRecording(event)
	ctx := context.Background()
	err = t.ConvCtrl.DB.SetConversationDone(ctx, state.RulesetID, callID)
	if err != nil {
		log.Printf("Error setting conversation done: %v", err)
		return // TODO: in its current state the conversation will be stuck in the DB
	}
	t.ConvCtrl.EndConversation(ctx, state, callID)
}

func (t *Telnyx) speakStartedProcedure(c *gin.Context, event Event) {
	log.Print("Speak started")
}

func (t *Telnyx) playbackStartedProcedure(c *gin.Context, event Event) {
	log.Print("playback started")
}

func (t *Telnyx) speakEndedProcedure(c *gin.Context, event Event) {
	state, err := decodeClientState(event.Data.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v. Ending call, possibly prematurely", err)
		t.EndCall(event.Data.Payload.CallControlID)
		return
	}
	if state.CurrentStep == state.TotalSteps-1 {
		log.Printf("Playback ended and conversation is done, ending call")
		t.EndCall(event.Data.Payload.CallControlID)
		return
	}
}

func (t *Telnyx) playbackEndedProcedure(c *gin.Context, event Event) {
	state, err := decodeClientState(event.Data.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v. Ending call, possibly prematurely", err)
		t.EndCall(event.Data.Payload.CallControlID)
		return
	}
	if state.CurrentStep == state.TotalSteps-1 {
		log.Printf("Playback ended and conversation is done, ending call")
		t.EndCall(event.Data.Payload.CallControlID)
		return
	}
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
		Url: url,
	})
}

func (t *Telnyx) recordingErrorProcedure(c *gin.Context, event Event) {
	log.Printf("Recording error: %v", event.Data.Payload.Reason)
}
