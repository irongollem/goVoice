package conversation

import (
	"context"
	"fmt"
	"goVoice/internal/email"
	"goVoice/internal/models"
	"goVoice/pkg/ai"
	"goVoice/pkg/audio"
	"goVoice/pkg/db"
	"goVoice/pkg/storage"
	"log"
	"sync"
	"time"
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
	AI       ai.AIProvider
	Email    *email.EmailProvider
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
		Responses: make(map[string]string),
	})

	// Grab the first step as conversation opener
	opener := ruleSet.Steps[0]
	clientState := models.ClientState{
		RulesetID:   rulesetID,
		CurrentStep: 0,
		Purpose:     opener.Purpose,
	}
	doneChan, errChan := c.broadcastNextStep(callID, &clientState, opener)

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
	validatedAnswer, err := c.validateAnswer(transcript, &rules.Steps[state.CurrentStep])
	if err != nil {
		log.Printf("Error validating answer, storing transcript: %v", err)
		c.storeTranscription(ctx, callID, state, rules, transcript)
	} else {
		log.Println("Validating succesful, storing validated answer")
		c.storeTranscription(ctx, callID, state, rules, validatedAnswer.Answer)
	}
	step, err := c.getResponse(rules, state, transcript)
	if err != nil {
		log.Printf("Error getting response for client: %v", err)
		// TODO tell the callee that something went wrong and handle gracefully
		return
	}

	nextState := models.ClientState{
		RulesetID:        state.RulesetID,
		CurrentStep:      state.CurrentStep + 1,
		Purpose:          rules.Steps[state.CurrentStep+1].Purpose,
		RecordingCount:   state.RecordingCount,
		RecordingPurpose: state.RecordingPurpose,
	}

	done, errChan := c.broadcastNextStep(callID, &nextState, step)
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

func (c *Controller) EndConversation(ctx context.Context, rulesetID string, callID string, recordingCount int) error {
	log.Printf("Ending conversation for %v", callID)
	conversation, err := c.DB.GetConversation(ctx, rulesetID, callID)
	if err != nil {
		log.Printf("Error getting responses from database: %v", err)
		return err
	}
	log.Println("Conversation collected, starting email process. Waiting for recordings to finish.")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)

	var recordings []models.Recording
Loop:
	for {
		select {
		case <-timeout:
			log.Printf("Timeout waiting for conversation to end, continuing incomplete for %v", callID)
			break Loop
		case <-ticker.C:
			recordings, err = c.DB.GetRecordings(ctx, rulesetID, callID)
			if err != nil {
				log.Printf("Error checking if conversation is complete: %v", err)
				return err
			}
			if recordings != nil && len(recordings) == recordingCount {
				break Loop
			}
		}
	}

	log.Println("Conversation complete, sending email.")

	ruleset, err := c.DB.GetRuleSet(ctx, rulesetID)
	if err != nil {
		log.Printf("Error getting ruleset from database: %v", err)
		return err
	}

	body := formatEmailBody(conversation.Responses, rulesetID, ruleset.Title, callID)
	var attachments [][]byte
	var attachmentsMutex sync.Mutex
	var attachmentNames []string
	var attachmentNamesMutex sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(recordings))

	log.Printf("Body formatting done; Downloading %v recordings", len(recordings))
	log.Println(recordings)

	errChan := make(chan error, len(recordings))

	for _, recording := range recordings {
		go func(rec *models.Recording) {
			defer wg.Done()

			recChan, recErrChan := c.Provider.GetRecordingMp3(rec)

			select {
			case recErr := <-recErrChan:
				if recErr != nil {
					log.Printf("Error getting recording: %v", recErr)
					errChan <- recErr
					return
				}
			case file := <-recChan:
				log.Printf("Got recording for %v added it to attachments", rec.Purpose)
				attachmentsMutex.Lock()
				attachmentNamesMutex.Lock()
				attachments = append(attachments, file)
				attachmentNames = append(attachmentNames, fmt.Sprintf("%s_%s-%s-%s.mp3", rulesetID, ruleset.Title, callID, rec.Purpose))
				attachmentsMutex.Unlock()
				attachmentNamesMutex.Unlock()
			}
		}(&recording)
	}
	log.Println("Waiting for recordings to finish downloading.")
	wg.Wait()

	close(errChan)
	if err, ok := <-errChan; ok {
		log.Printf("Error getting recording: %v", err)
		return err
	}
	log.Println("Recordings downloaded, sending email.")
	err = c.Email.SendEmailWithAttachment(ctx, ruleset.Client.Email, ruleset.Title, body, attachments, attachmentNames)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	err = c.DB.DeleteConversation(ctx, rulesetID, callID)
	if err != nil {
		log.Printf("Error deleting conversation from database: %v", err)
		return err
	}

	return nil
}

func (c *Controller) ProcessRecording(ctx context.Context, rulesetId string, callID string, recording *models.Recording) error {
	// Set the recording on the conversation
	err := c.DB.SetRecording(ctx, rulesetId, callID, recording)
	if err != nil {
		log.Printf("Error setting recordings on conversation: %v", err)
		return err
	}
	return nil
}
