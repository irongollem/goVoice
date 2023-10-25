package conversation

import (
	"encoding/json"
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
// rules and combines those with the transcription to send to the
// LLM controller to determine a response. The response is then sent to the audio
// processor to be converted to audio and sent back to the caller.
type Controller struct {
	Provider audio.CallProvider
	Storage  storage.StorageProvider
	DB       db.DbProvider
}

func (c *Controller) StartConversation(callId string) {
	conv, err := c.getRules(callId)
	if err != nil {
		log.Printf("Error getting conversation rules: %v", err)
		return
	}

	// Grab the first step as conversation opener
	opener := conv.Steps[0]
	clientState := ClientStateToJSON(&opener, 0)

	doneChan, errChan := c.Provider.SpeakText(callId, opener.Text, clientState)
	select {
	case <-doneChan:
		log.Printf("Successfully sent conversation opener to caller")
	case err := <-errChan:
		log.Printf("Error sending conversation opener to caller: %v", err)
		// TODO: do we abandon the call?
	}
}

func (c *Controller) ProcessTranscription(callId string, transcript string, clientState string) {
	state := JSONToClientState(&clientState)
	rules, err := c.getRules(callId)

	// in case people are still talking after the conversation is over
	wasFinalStep := len(rules.Steps) == state.Index
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

	storeTranscription(callId, state.Purpose, transcript)

	// combine transcription with conversation rules
	step, err := c.getResponse(rules, state, transcript)
	if err != nil {
		log.Printf("Error getting response for client: %v", err)
		// TODO tell the callee that something went wrong and handle gracefully
		return
	}

	nextClientState := ClientStateToJSON(&step, state.Index+1)
	done, errChan := c.Provider.SpeakText(callId, step.Text, nextClientState)
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

func (c *Controller) getRules(callId string) (models.ConversationRuleSet, error) {
	//  after pilot we should actually fetch this
	file, err := os.Open("pilot.json")
	if err != nil {
		log.Printf("Error opening pilot.json file: %v", err)
		return models.ConversationRuleSet{}, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading pilot.json file: %v", err)
		return models.ConversationRuleSet{}, err
	}

	var ruleSet models.ConversationRuleSet
	err = json.Unmarshal(byteValue, &ruleSet)
	if err != nil {
		log.Printf("Error unmarshalling pilot.json file: %v", err)
		return models.ConversationRuleSet{}, err
	}

	return ruleSet, nil
}

func (c *Controller) getResponse(rules models.ConversationRuleSet, state *ClientState, transcript string) (models.ConversationStep, error) {
	if rules.Simple {
		return getSimpleResponse(rules, state), nil
	} else {
		return getAdvancedResponse(rules, state, transcript)
	}
}

// get a response using an LLM
func getAdvancedResponse(rules models.ConversationRuleSet, state *ClientState, transcript string) (models.ConversationStep, error) {
	panic("unimplemented")
}

// get a response using a simple call script
func getSimpleResponse(rules models.ConversationRuleSet, state *ClientState) models.ConversationStep {
	return rules.Steps[state.Index+1]
}
