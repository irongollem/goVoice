package conversation

import (
	"encoding/base64"
	"encoding/json"
	"log"
)

// JSONToClientState decodes the given client state string and returns a ConversationStep pointer.
// If the client state is nil or decoding fails, it returns nil.
func JSONToClientState(clientState *string) *ClientState {
	if clientState == nil {
		return nil
	}

	stateJSON, err := base64.StdEncoding.DecodeString(*clientState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
		return nil
	}

	var step ClientState
	err = json.Unmarshal(stateJSON, &step)
	if err != nil {
		log.Printf("Error unmarshalling client state: %v", err)
		return nil
	}

	return &step
}

// ClientStateToJSON encodes the given ConversationStep pointer and returns a base64-encoded string pointer.
// If the step is nil or encoding fails, it returns nil.
func ClientStateToJSON(step *ConversationStep, stepIndex int) *string {
	if step == nil {
		return nil
	}

	state := ClientState{
		Index:   stepIndex,
		Purpose: step.Purpose,
	}

	stepJSON, err := json.Marshal(state)
	if err != nil {
		log.Printf("Error marshalling client state: %v", err)
		return nil
	}

	encodedState := base64.StdEncoding.EncodeToString(stepJSON)
	return &encodedState
}

func storeTranscription(callId string, purpose string, transcript string) {

}
