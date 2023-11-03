package conversation

import (
	"context"
	"encoding/json"
	"goVoice/internal/models"
	"goVoice/pkg/audio"
	"goVoice/pkg/db"
	"goVoice/pkg/email"
	"goVoice/pkg/storage"
	"io"
	"log"
	"os"
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

func (c *Controller) StartConversation(callId string) {
	conv, err := c.getRules(callId)
	if err != nil {
		log.Printf("Error getting conversation rules: %v", err)
		return
	}

	// Grab the first step as conversation opener
	opener := conv.Steps[0]
	clientState := &models.ClientState{
		RulesetID: callId,
		CurrentStep: 0,
	}

	doneChan, errChan := c.Provider.SpeakText(callId, opener.Text, clientState)
	select {
	case <-doneChan:
		log.Printf("Successfully sent conversation opener to caller")
	case err := <-errChan:
		log.Printf("Error sending conversation opener to caller: %v", err)
		// TODO: do we abandon the call?
	}
}

func (c *Controller) ProcessTranscription(ctx context.Context, callId string, transcript string, state *models.ClientState) {
	rules, err := c.getRules(callId)
	if err != nil {
		log.Printf("Error getting conversation rules: %v", err)
		return
	}

	// in case people are still talking after the conversation is over
	wasFinalStep := len(rules.Steps) == state.CurrentStep - 1
	if wasFinalStep {
		done, errChan := c.Provider.EndCall(callId)
		select {
		case <-done:
			return
		case err := <-errChan:
			log.Printf("Error ending call: %v", err)
			return
		}
	}

	c.storeTranscription(ctx, callId, state, rules,transcript)

	// combine transcription with conversation rules
	step, err := c.getResponse(rules, state, transcript)
	if err != nil {
		log.Printf("Error getting response for client: %v", err)
		// TODO tell the callee that something went wrong and handle gracefully
		return
	}

	nextState := models.ClientState{
		RulesetID: callId,
		CurrentStep: state.CurrentStep + 1,
	}
	done, errChan := c.Provider.SpeakText(callId, step.Text, &nextState)
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

func (c *Controller) getRules(callId string) (*models.ConversationRuleSet, error) {
	//  after pilot we should actually fetch this
	file, err := os.Open("pilot.json")
	if err != nil {
		log.Printf("Error opening pilot.json file: %v", err)
		return &models.ConversationRuleSet{}, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading pilot.json file: %v", err)
		return &models.ConversationRuleSet{}, err
	}

	var ruleSet *models.ConversationRuleSet
	err = json.Unmarshal(byteValue, &ruleSet)
	if err != nil {
		log.Printf("Error unmarshalling pilot.json file: %v", err)
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
	return rules.Steps[state.CurrentStep+1]
}

func (c *Controller) EndConversation(ctx context.Context, rulesetId string, callId string) error {
	// Get the stored responses from the database
	responses, err := c.DB.GetResponses(ctx, rulesetId, callId)
	if err != nil {
		log.Printf("Error getting responses from database: %v", err)
		return err
	}

	// Get the recording or a link to it from the storage provider
	recording, err := c.Storage.GetRecording(ctx, rulesetId, callId)
	if err != nil {
		log.Printf("Error getting recording from storage: %v", err)
		return err
	}

	// send the data to sendgrid to be emailed to the client
	err = c.email.SendConversationEmail(ctx, recording, responses)

	panic("unimplemented")
}
