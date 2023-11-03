package telnyx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goVoice/internal/models"
	"log"
	"math"
	"net/http"
	"time"
)

func (t *Telnyx) sendCommand(callControlId string, command string, payload *CommandPayload) (*http.Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return nil, err
	}

	var resp *http.Response
	var retries int

	for retries = 0; retries < 5; retries++ {
		req, err := http.NewRequest("POST", t.CommandPath+callControlId+"/actions/"+command, bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+t.APIKey)

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error sending command: %v", err)
			continue
		}

		if resp.StatusCode >= 500 && resp.StatusCode < 600 {
			log.Printf("Received %d status code, retrying", resp.StatusCode)
			time.Sleep(time.Duration(math.Pow(2, float64(retries))) * time.Second)
			continue
		}

		break
	}

	if retries == 5 {
		return nil, fmt.Errorf("failed to send command after %d retries", retries)
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
		_, err := t.sendCommand(event.Payload.CallControlID, "answer", &answerPayload)
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
		_, err := t.sendCommand(event.Payload.CallControlID, "transcription_start", &transcriptionPayload)
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
		_, err := t.sendCommand(event.Payload.CallControlID, "record_start", &recordingPayload)
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
		_, err := t.sendCommand(callId, "speak_start", &speakPayload)
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
		_, err := t.sendCommand(callId, "hangup", nil)
		if err != nil {
			log.Printf("Error ending call: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}
