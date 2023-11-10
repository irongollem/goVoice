package telnyx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goVoice/internal/models"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/avast/retry-go"
)
func (t *Telnyx) sendCommandToCallCommandsAPI (callControlId string, command string, payload *CommandPayload) (*http.Response, error) {
	return t.sendCommand("POST", payload, "calls", callControlId, "actions",  command)
}

func (t *Telnyx) sendCommand( httpMethod string, payload interface{}, pathParams ...string) (*http.Response, error) {
	var payloadBytes []byte
	var err error
	
	switch v := payload.(type) {
		case *CommandPayload:
			payloadBytes, err = json.Marshal(v)
		case *CredentialsPayload:
			payloadBytes, err = json.Marshal(v)
		default:
			log.Printf("Unknown payload type: %T", v)
			return nil, fmt.Errorf("unknown payload type: %T", v)
	}
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return nil, err
	}

	for _, param := range pathParams {
		newPath := path.Join(t.APIUrl.Path, param)
		t.APIUrl.Path = newPath
	}

	var resp *http.Response
	if err = retry.Do(
		func() error{
			req, err := http.NewRequest(httpMethod, t.APIUrl.String(), bytes.NewBuffer(payloadBytes))
			if err != nil {
				log.Printf("Error creating request: %v", err)
				return err
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+t.APIKey)

			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("Error sending command: %v", err)
				return err
			}

			if resp.StatusCode >= 500 && resp.StatusCode < 600 {
				log.Printf("Received %d status code, retrying", resp.StatusCode)
				return fmt.Errorf("received %d status code, retrying", resp.StatusCode)
			}

			return nil
		},
		retry.Attempts(3),
		retry.Delay(time.Second),
	); err != nil {
		log.Printf("Error sending command: %v", err)
		return nil, err
	}


	return resp, nil
}

func (t *Telnyx) answerCall(event Event) (chan bool, chan error) {
	done := make(chan bool)
	errChan := make(chan error, 1)

	answerPayload := CommandPayload{
		SendSilenceWhenIdle: true,
		ClientState:         event.Payload.ClientState,
	}

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(event.Payload.CallControlID, "answer", &answerPayload)
		if err != nil {
			log.Printf("Error answering call: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) startTranscription(event Event) (chan bool, chan error) {
	done := make(chan bool)
	errChan := make(chan error, 1)

	transcriptionPayload := CommandPayload{
		ClientState:         event.Payload.ClientState,
		CommandId:           "transcription-1", // TODO make sure this if fine
		Language:            "nl",              // TODO: try using auto_detect and talk other languages
		TranscriptionEngine: "B",               // A is google, B is telnyx
	}

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(event.Payload.CallControlID, "transcription_start", &transcriptionPayload)
		if err != nil {
			log.Printf("Error starting transcription: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) startRecording(event Event) (chan bool, chan error) {
	done := make(chan bool)
	errChan := make(chan error, 1)

	recordingPayload := CommandPayload{
		ClientState: event.Payload.ClientState,
		CommandId:   "recording-1", // TODO make sure this if fine
		Format:      "mp3",
		Channels:    "single",
		Trim:        "trim-silence",
	}

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(event.Payload.CallControlID, "record_start", &recordingPayload)
		if err != nil {
			log.Printf("Error starting recording: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) SpeakText(callId string, text string, clientState *models.ClientState) (chan bool, chan error) {
	done := make(chan bool)
	errChan := make(chan error, 1)

	state, _ := encodeClientState(clientState)

	speakPayload := CommandPayload{
		CommandId:   "speak-1", // TODO make sure this if fine
		Language:    "nl-NL",
		Voice:       "male",
		Payload:     text,
		ClientState: state,
	}

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(callId, "speak_start", &speakPayload)
		if err != nil {
			log.Printf("Error starting speak: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) EndCall(callId string) (chan bool, chan error) {
	done := make(chan bool)
	errChan := make(chan error, 1)

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(callId, "hangup", nil)
		if err != nil {
			log.Printf("Error ending call: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) GetRecordings(callId string) (chan []Recording, chan error) {
	done := make(chan []Recording)
	errChan := make(chan error, 1)

	go func() {
		res, err := t.sendCommand(callId, nil, "GET", "recordings")
		if err != nil {
			log.Printf("Error getting recordings: %v", err)
			errChan <- err
			return
		}
		defer res.Body.Close()

		var recordings []Recording
		if err := json.NewDecoder(res.Body).Decode(&recordings); err != nil {
			log.Printf("Error decoding recordings: %v", err)
			errChan <- err
			return
		}

		var matchingRecordings []Recording
		for _, r := range recordings {
			if r.CallControlID == callId {
				matchingRecordings = append(matchingRecordings, r)
			}
		}
		if len(matchingRecordings) == 0 {
			errChan <- fmt.Errorf("no recording found for call %s", callId)
			return
		}
		done <- matchingRecordings
	}()

	return done, errChan
}
