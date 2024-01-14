package telnyx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goVoice/internal/models"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
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
	_, isPlayAudioPayload := payload.(*PlayAudio)
	_, isNoiseSuppressionPayload := payload.(*NoiseSuppressionPayload)
	_, isTranscriptionPayload := payload.(*TranscriptionPayload)

	if payload != nil && !isSimplePayload && !isAnswerPayload && !isUpdateClientStatePayload && !isGatherPayload && !isRecordStartPayload && !isSpeakTextPayload && !isPlayAudioPayload && !isNoiseSuppressionPayload && !isTranscriptionPayload {
		log.Printf("Unknown payload type: %T", payload)
		return nil, fmt.Errorf("unknown payload type: %T", payload)
	}

	newPath := ""
	for _, param := range pathParams {
		newPath = path.Join(newPath, param)
	}
	newURL := *t.APIUrl
	newURL.Path = path.Join(newURL.Path, newPath)

	var req *http.Request
	var err error
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshaling payload: %v", err)
			return nil, err
		}
		req, err = http.NewRequest(httpMethod, newURL.String(), bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(httpMethod, newURL.String(), nil)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return nil, err
		}
	}

	var resp *http.Response
	if err = retry.Do(
		func() error {
			req.Header.Set("Authorization", "Bearer "+t.APIKey)

			dump, err := httputil.DumpRequest(req, true)
			if err != nil {
				log.Printf("Error dumping request: %v", err)
			} else {
				log.Printf("Request: %s", dump)
			}

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
	state, err := encodeClientState(&models.ClientState{
		RulesetID:   "LeeuwardenPilot", // FIXME this should be dynamic
		CurrentStep: 0,
	})
	if err != nil {
		log.Printf("Error encoding client state: %v", err)
		errChan <- err
		return done, errChan
	}

	answerPayload := &AnswerPayload{
		SendSilenceWhenIdle: true,
		ClientState:         state,
		CommandID:           generateCommandID(event.Data.Payload.CallControlID, "answer", state),
	}

	go func() {
		res, err := t.sendCommandToCallCommandsAPI(event.Data.Payload.CallControlID, "answer", answerPayload)
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
		Language:            "nl",
		TranscriptionEngine: "A",  // A is google, B is telnyx
	}

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(event.Data.Payload.CallControlID, "transcription_start", transcriptionPayload)
		if err != nil {
			log.Printf("Error starting transcription: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) stopTranscription(event Event) (chan bool, chan error) {
	log.Printf("Stopping transcription for call %s", event.Data.Payload.CallControlID)
	done := make(chan bool)
	errChan := make(chan error, 1)

	payload := &SimplePayload{
		ClientState: event.Data.Payload.ClientState,
		CommandID:   generateCommandID(event.Data.Payload.CallControlID, "transcription_stop", event.Data.Payload.ClientState),
	}
	
	go func() {
		_, err := t.sendCommandToCallCommandsAPI(event.Data.Payload.CallControlID, "transcription_stop", payload)
		if err != nil {
			log.Printf("Error stopping transcription: %v", err)
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

	var newState string
	state, err := decodeClientState(event.Data.Payload.ClientState)
	if err != nil {
		log.Printf("Error decoding client state: %v", err)
		newState = event.Data.Payload.ClientState
	} else {
		state.RecordingCount++
		newState, err = encodeClientState(state)
		if err != nil {
			log.Printf("Error encoding client state: %v", err)
			newState = event.Data.Payload.ClientState
		}
	}

	recordingPayload := &RecordStartPayload{
		ClientState: newState,
		CommandID:   generateCommandID(event.Data.Payload.CallControlID, "record_start", event.Data.Payload.ClientState),
		Format:      "mp3",
		Channels:    "single",
		Trim:        "trim-silence",
	}

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(event.Data.Payload.CallControlID, "record_start", recordingPayload)
		if err != nil {
			log.Printf("Error starting recording: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) stopRecording(event Event) (chan bool, chan error) {
	log.Printf("Stopping recording for call %s", event.Data.Payload.CallControlID)
	done := make(chan bool)
	errChan := make(chan error, 1)

	payload := &SimplePayload{
		ClientState: event.Data.Payload.ClientState,
		CommandID:   generateCommandID(event.Data.Payload.CallControlID, "record_stop", event.Data.Payload.ClientState),
	}

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(event.Data.Payload.CallControlID, "record_stop", payload)
		if err != nil {
			log.Printf("Error stopping recording: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) SpeakText(CallControlID string, text string, clientState *models.ClientState) (chan bool, chan error) {
	log.Printf("Speaking text for call %s", CallControlID)
	command := "speak"
	done := make(chan bool)
	errChan := make(chan error, 1)

	state, _ := encodeClientState(clientState)

	speakPayload := &SpeakTextPayload{
		CommandID:   generateCommandID(CallControlID, command, state),
		Language:    "nl-NL",
		Voice:       "male",
		Payload:     text,
		ClientState: state,
	}

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(CallControlID, command, speakPayload)
		if err != nil {
			log.Printf("Error starting speak: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) EndCall(callID string) (chan bool, chan error) {
	log.Printf("Ending call %s", callID)
	done := make(chan bool)
	errChan := make(chan error, 1)

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(callID, "hangup", nil)
		if err != nil {
			log.Printf("Error ending call: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}

func (t *Telnyx) GetRecording(recordingId string) (chan *Recording, chan error) {
	log.Printf("Getting recording: %s", recordingId)
	done := make(chan *Recording)
	errChan := make(chan error, 1)

	go func() {
		res, err := t.sendCommand("GET", nil, "recordings", recordingId)
		if err != nil {
			log.Printf("Error getting recordings: %v", err)
			errChan <- err
			return
		}
		defer res.Body.Close()

		log.Println("-------- Response Body Start --------")
		// Read and log the response body
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			errChan <- err
			return
		}
		log.Println(string(body))
		log.Println("-------- Response Body End --------")

		var recordingResponse RecordingResponse
		if err := json.NewDecoder(res.Body).Decode(&recordingResponse); err != nil {
			log.Printf("Error decoding recordings: %v", err)
			errChan <- err
			return
		}

		done <- recordingResponse.Data
	}()

	return done, errChan
}

func (t *Telnyx) GetRecordingMp3(recording *models.Recording) (chan []byte, chan error) {
	log.Printf("Getting recording mp3: %s", recording.Purpose)
	done := make(chan []byte)
	errChan := make(chan error, 1)

	go func() {
		req, err := http.NewRequest("GET", recording.Url, nil)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			errChan <- err
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error getting recording mp3: %v", err)
			errChan <- err
			return
		}
		defer res.Body.Close()
		log.Printf("Response status: %s, downloaded the file, about to read body", res.Status)
		mp3Bytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error reading recording mp3: %v", err)
			errChan <- err
			return
		}
		log.Printf("Read the body, about to return")

		done <- mp3Bytes
	}()

	return done, errChan
}

func (t *Telnyx) PlayAudioUrl(callControlID string, url string, clientState *models.ClientState)(chan bool, chan error) {
	log.Printf("Playing audio url: %s", url)
	done := make(chan bool)
	errChan := make(chan error, 1)

	state, _ := encodeClientState(clientState)

	payload := &PlayAudio{
		CommandID:   generateCommandID(callControlID, "playAudio", state),
		AudioUrl: url,
		ClientState: state,
	}

	go func() {
		_, err := t.sendCommandToCallCommandsAPI(callControlID, "playback_start", payload)
		if err != nil {
			log.Printf("Error playing audio url: %v", err)
			errChan <- err
		}
		done <- true
	}()

	return done, errChan
}
