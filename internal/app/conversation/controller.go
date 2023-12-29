package conversation

import (
	"context"
	"fmt"
	"goVoice/internal/email"
	"goVoice/internal/models"
	"goVoice/pkg/audio"
	"goVoice/pkg/db"
	"goVoice/pkg/storage"
	"log"
)

// The conversation controller orchestrates the incoming audio, sends it
// to the audio processor to be transcribed, then fetches the conversation
// rules and combines thapose with the transcription to send to the
// LLM controller to determine a response. The response is then sent to the audio
// processor to be converted to audio and sent back to the caller.
type Controller struct {
	Provider audio.CallProvider
	Storage  storage.StorageProvider
	DB       db.DbProvider
	email    email.EmailProvider
}

func (c *Controller) StartConversation(rulesetID string, callID string) {
	ruleSet, err := c.getRules(rulesetID)
	if err != nil {
		log.Printf("Error getting conversation rules: %v", err)
		return
	}

	c.DB.AddConversation(context.Background(), rulesetID, &models.Conversation{
		ID:        callID,
		RulesetID: rulesetID,
		Responses: []models.ConversationStepResponse{},
	})

	// Grab the first step as conversation opener
	opener := ruleSet.Steps[0]
	clientState := &models.ClientState{
		RulesetID:   rulesetID,
		CurrentStep: 0,
	}

	doneChan, errChan := c.Provider.SpeakText(callID, opener.Text, clientState)
	select {
	case <-doneChan:
		log.Printf("Successfully sent conversation opener to caller")
	case err := <-errChan:
		log.Printf("Error sending conversation opener to caller: %v", err)
		// TODO: do we abandon the call?
	}
}

func (c *Controller) ProcessTranscription(ctx context.Context, callID string, transcript string, state *models.ClientState) {
	rules, err := c.getRules(state.RulesetID)
	if err != nil {
		log.Printf("Error getting conversation rules: %v", err)
		return
	}

	// in case people are still talking after the conversation is over
	wasFinalStep := len(rules.Steps) == state.CurrentStep-1
	if wasFinalStep {
		done, errChan := c.Provider.EndCall(callID)
		select {
		case <-done:
			return
		case err := <-errChan:
			log.Printf("Error ending call: %v", err)
			return
		}
	}

	c.storeTranscription(ctx, callID, state, rules, transcript)

	step, err := c.getResponse(rules, state, transcript)
	if err != nil {
		log.Printf("Error getting response for client: %v", err)
		// TODO tell the callee that something went wrong and handle gracefully
		return
	}

	nextState := models.ClientState{
		RulesetID:   state.RulesetID,
		CurrentStep: state.CurrentStep + 1,
	}
	done, errChan := c.Provider.SpeakText(callID, step.Text, &nextState)
	select {
	case <-done:
		return
	case err := <-errChan:
		log.Printf("Error sending response to caller: %v", err)
		// TODO: do we abandon the call?
		// and what do we do with the state and recording?
		// do we save it to the database and send it to the client?
	}
}

func (c *Controller) getRules(rulesetID string) (*models.ConversationRuleSet, error) {
	context := context.Background()
	ruleSet, err := c.DB.GetRuleSet(context, rulesetID)
	if err != nil {
		log.Printf("Error fetching ruleset from DB: %v", err)
		return &models.ConversationRuleSet{}, err
	}

	return ruleSet, nil
}

func (c *Controller) getResponse(rules *models.ConversationRuleSet, state *models.ClientState, transcript string) (models.ConversationStep, error) {
	if rules.Simple {
		return getSimpleResponse(rules, state), nil
	} else {
		return getAdvancedResponse(rules, state, transcript)
	}
}

// get a response using an LLM
func getAdvancedResponse(rules *models.ConversationRuleSet, state *models.ClientState, transcript string) (models.ConversationStep, error) {
	panic("unimplemented")
}

// get a response using a simple call script
func getSimpleResponse(rules *models.ConversationRuleSet, state *models.ClientState) models.ConversationStep {
	if len(rules.Steps) >= state.CurrentStep+1 {
		return rules.Steps[state.CurrentStep+1]
	} else {
		return models.ConversationStep{
			Text: "Bedankt voor het bellen, tot ziens!",
		}
	}
}

func (c *Controller) EndConversation(ctx context.Context, rulesetID string, callID string) error {
	recording, err := c.DB.IsConversationComplete(ctx, rulesetID, callID)
	if err != nil {
		log.Printf("Error checking if conversation is complete: %v", err)
		return err
	}
	if recording == nil {
		log.Print("Conversation is not complete")
		return nil
	}

	ruleset, err := c.DB.GetRuleSet(ctx, rulesetID)
	if err != nil {
		log.Printf("Error getting ruleset from database: %v", err)
		return err
	}

	body, _ := c.retrieveResponsesAsTable(ctx, ruleset, callID)

	recChan, errChan := c.Provider.GetRecordingMp3(recording)
	select {
	case recording := <-recChan:
		filename := fmt.Sprintf("%s_%s-%s.mp3", rulesetID, ruleset.Title, callID)

		err = c.email.SendEmailWithAttachment(ctx, ruleset.Client.Email, ruleset.Title, *body, recording, filename)
		if err != nil {
			log.Printf("Error sending email: %v", err)
			return err
		}
		return nil
	case err := <-errChan:
		log.Printf("Error getting recording: %v", err)
		return err
	}
	// TODO: delete the recording from storage and conversation from db if data is send off successfully
}

func (c *Controller) ProcessRecording(ctx context.Context, rulesetId string, callID string, recording *models.Recording) error {
	// Set the recording on the conversation
	err := c.DB.SetRecordings(ctx, rulesetId, callID, recording)
	if err != nil {
		log.Printf("Error setting recordings on conversation: %v", err)
		return err
	}
	// If the conversation is not complete, EndConversation will deal with that
	return c.EndConversation(ctx, rulesetId, callID)
}
