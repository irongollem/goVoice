package telnyx

import (
	"encoding/base64"
	"encoding/json"
	"goVoice/internal/models"
	"log"

	"github.com/google/uuid"
)

func decodeClientState (encState string) (*models.ClientState, error) {
	var state models.ClientState
	decodeState, err := base64.StdEncoding.DecodeString(encState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
		return nil, err
	}
	if err = json.Unmarshal(decodeState, &state); err != nil {
		log.Printf("Error unmarshalling client state: %v", err)
		return nil, err
	}
	return &state, nil
}

func encodeClientState (state *models.ClientState) (string, error) {
	encodeState, err := json.Marshal(state)
	if err != nil {
		log.Printf("Error marshalling client state: %v", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encodeState), nil
}

func convertToRecording (recordings []Recording) ([]models.Recording) {
	var genericRecordings []models.Recording
	for _, recording := range recordings {
		genericRecordings = append(genericRecordings, models.Recording{
			Url: recording.DownloadUrls.Mp3,
			ID: recording.ID,
			ConversationID: recording.CallControlID,
		})
	}
	return genericRecordings
}

func generateCommandID (CallControlID, funcName string, clientState string) string {
	idString := CallControlID + funcName + clientState + "v1"

	id := uuid.NewMD5(uuid.Nil, []byte(idString))
	return id.String()
}