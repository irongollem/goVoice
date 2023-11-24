package conversation

import (
	"context"
	"encoding/json"
	"fmt"
	"goVoice/internal/email"
	"goVoice/internal/models"
	"goVoice/pkg/audio"
	"goVoice/pkg/db"
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
		RulesetID:   callId,
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
	wasFinalStep := len(rules.Steps) == state.CurrentStep-1
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

	c.storeTranscription(ctx, callId, state, rules, transcript)

	step, err := c.getResponse(rules, state, transcript)
	if err != nil {
		log.Printf("Error getting response for client: %v", err)
		// TODO tell the callee that something went wrong and handle gracefully
		return
	}

	nextState := models.ClientState{
		RulesetID:   callId,
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
	file, err := os.Open("./pilot.json")
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
	if len(rules.Steps) >= state.CurrentStep+1 {
		return rules.Steps[state.CurrentStep+1]
	} else {
		return models.ConversationStep{
			Text: "Bedankt voor het bellen, tot ziens!",
		}
	}
}

// FIXME get recording from storage or through TELNYX is unclear, so far to the DB we wrote the ID and url as provided by telnyx. If that also defines the file ID we could use our direct connection to storage, otherwise we need to use the telnyx API to get the recording
func (c *Controller) EndConversation(ctx context.Context, rulesetId string, callId string) error {
	// Typically this is called when the caller hangs up, or when the conversation is complete
	// Which is triggered when recording is done OR when the LLM determines the conversation is complete
	onOwnBucket := false
	// Get the recording or a link to it from the storage provider
	var err error
	var reader io.ReadCloser
	if onOwnBucket {
		reader, err = c.Storage.GetRecording(ctx, rulesetId, callId)
		if err != nil {
			log.Printf("Error getting recording from storage: %v", err)
			return err
		}
	} else {
		// FIXME get recording from telnyx
		// reader, err = c.Provider.GetRecording(ctx, rulesetId, callId)
		// if err != nil {
		// 	log.Printf("Error getting recording from telnyx: %v", err)
		// 	return err
		// }
	}
	defer reader.Close()

	recording, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("Error reading recording from storage: %v", err)
		return err
	}
	filename := fmt.Sprintf("recording-%s-%s.mp3", rulesetId, callId)

	ruleset, err := c.DB.GetRuleSet(ctx, rulesetId)
	if err != nil {
		log.Printf("Error getting ruleset from database: %v", err)
		return err
	}

	body, _ := c.retrieveResponsesAsTable(ctx, ruleset, callId)
	// if there is an error in the responses, we still want to send the recording

	err = c.email.SendEmailWithAttachment(ctx, ruleset.Client.Email, ruleset.Title, *body, recording, filename)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	return nil
}

func (c *Controller) ProcessRecording(ctx context.Context, rulesetId string, callId string, recording *models.Recording) error {
	// Set the recording on the conversation
	err := c.DB.SetRecordings(ctx, rulesetId, callId, recording)
	if err != nil {
		log.Printf("Error setting recordings on conversation: %v", err)
		return err
	}
	// If the conversation is not complete, EndConversation will deal with that
	return c.EndConversation(ctx, rulesetId, callId)
}
