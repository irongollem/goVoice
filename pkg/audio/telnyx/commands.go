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

func (t *Telnyx) sendCommandToCallCommandsAPI(callControlId string, command string, payload interface{}) (*http.Response, error) {
	return t.sendCommand("POST", payload, "calls", callControlId, "actions", command)
}

func (t *Telnyx) sendCommand(httpMethod string, payload interface{}, pathParams ...string) (*http.Response, error) {
	_, isSimplePayload := payload.(*SimplePayload)
	_, isAnswerPayload := payload.(*AnswerPayload)
	_, isUpdateClientStatePayload := payload.(*UpdateClientStatePayload)
	_, isGatherPayload := payload.(*GatherPayload)
	_, isRecordStartPayload := payload.(*RecordStartPayload)
	_, isSpeakTextPayload := payload.(*SpeakTextPayload)
	_, isNoiseSuppressionPayload := payload.(*NoiseSuppressionPayload)
	_, isTranscriptionPayload := payload.(*TranscriptionPayload)

	if !isSimplePayload && !isAnswerPayload && !isUpdateClientStatePayload && !isGatherPayload && !isRecordStartPayload && !isSpeakTextPayload && !isNoiseSuppressionPayload && !isTranscriptionPayload {
		log.Printf("Unknown payload type: %T", payload)
		return nil, fmt.Errorf("unknown payload type: %T", payload)
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return nil, err
	}

	for _, param := range pathParams {
		newPath := path.Join(t.APIUrl.Path, param)
		t.APIUrl.Path = newPath
	}
	log.Printf("Sending command to %s: %s; With api key %s", t.APIUrl.String(), string(payloadBytes), t.APIKey)
	var resp *http.Response
	if err = retry.Do(
		func() error {
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
	log.Printf("Answering call %s", event.Data.Payload.CallControlID)
	done := make(chan bool)
	errChan := make(chan error, 1)

	answerPayload := &AnswerPayload{
		SendSilenceWhenIdle: true,
		ClientState:         event.Data.Payload.ClientState,
		CommandID:           generateCommandID(event.Data.Payload.CallControlID, "answer", event.Data.Payload.ClientState),
	}

	go func() {
		res, err := t.sendCommandToCallCommandsAPI(event.Data.Payload.CallLegID, "answer", &answerPayload)
		if err != nil {
			log.Printf("Error answering call: %v", err)
			errChan <- err
		}
		if res.StatusCode >= http.StatusBadRequest {
			log.Printf("Error answering call: %v", res.Status)
			errChan <- fmt.Errorf("error answering call: %v", res.Status)
		}
		log.Printf("Answered call %s", event.Data.Payload.CallControlID)
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) startTranscription(event Event) (chan bool, chan error) {
	log.Printf("Starting transcription for call %s", event.Data.Payload.CallControlID)
	done := make(chan bool)
	errChan := make(chan error, 1)

	transcriptionPayload := &TranscriptionPayload{
		ClientState:         event.Data.Payload.ClientState,
		CommandID:           generateCommandID(event.Data.Payload.CallControlID, "transcription_start", event.Data.Payload.ClientState),
		Language:            "nl",              // TODO: try using auto_detect and talk other languages
		TranscriptionEngine: "B",               // A is google, B is telnyx
	}

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(event.Data.Payload.CallControlID, "transcription_start", &transcriptionPayload)
		if err != nil {
			log.Printf("Error starting transcription: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) startRecording(event Event) (chan bool, chan error) {
	log.Printf("Starting recording for call %s", event.Data.Payload.CallControlID)
	done := make(chan bool)
	errChan := make(chan error, 1)

	recordingPayload := RecordStartPayload{
		ClientState: event.Data.Payload.ClientState,
		CommandID:   generateCommandID(event.Data.Payload.CallControlID, "record_start", event.Data.Payload.ClientState),
		Format:      "mp3",
		Channels:    "single",
		Trim:        "trim-silence",
	}

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(event.Data.Payload.CallControlID, "record_start", &recordingPayload)
		if err != nil {
			log.Printf("Error starting recording: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) SpeakText(CallControlID string, text string, clientState *models.ClientState) (chan bool, chan error) {
	log.Printf("Speaking text for call %s", CallControlID)
	done := make(chan bool)
	errChan := make(chan error, 1)

	state, _ := encodeClientState(clientState)

	speakPayload := SpeakTextPayload{
		CommandID:   generateCommandID(CallControlID, "speak_start", state),
		Language:    "nl-NL",
		Voice:       "male",
		Payload:     text,
		ClientState: state,
	}

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(CallControlID, "speak_start", &speakPayload)
		if err != nil {
			log.Printf("Error starting speak: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) EndCall(callId string) (chan bool, chan error) {
	log.Printf("Ending call %s", callId)
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
	log.Printf("Getting recordings for call %s", callId)
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
