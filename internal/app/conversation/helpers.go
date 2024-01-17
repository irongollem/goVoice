package conversation

import (
	"context"
	"encoding/json"
	"fmt"
	"goVoice/internal/models"
	"goVoice/pkg/ai"
	"log"
	"strings"
)

func (c *Controller) storeTranscription(ctx context.Context, callID string, state *models.ClientState, ruleSet *models.ConversationRuleSet, transcript string) {
	log.Printf("Storing transcription: %s", transcript)
	c.DB.AddResponse(ctx, state.RulesetID, callID, &models.ConversationStepResponse{
		Purpose:  ruleSet.Steps[state.CurrentStep].Purpose,
		Response: transcript,
	})
}

func formatEmailBody(responses map[string]string, rulesetID string, rulesetTitle string, callID string) string {
	var sb strings.Builder
	sb.WriteString("<table style='width: 100%; border-collapse: collapse;'>\n")
	sb.WriteString(fmt.Sprintf("<tr style='background-color: #f2f2f2;'><td style='border: 1px solid #ddd; padding: 8px;'>CallID</td><td style='border: 1px solid #ddd; padding: 8px;'>%s</td></tr>\n", callID))
	i := 0
	for purpose, answer := range responses {
			color := "#f2f2f2"
			if i%2 == 0 {
					color = "#ddf"
			}
			sb.WriteString(fmt.Sprintf("<tr style='background-color: %s;'><td style='border: 1px solid #ddd; padding: 8px;'>%s</td><td style='border: 1px solid #ddd; padding: 8px;'>%s</td></tr>\n", color, purpose, answer))
			i++
	}
	sb.WriteString("</table>")
	return sb.String()
}

func (c *Controller) validateAnswer(answer string, step *models.ConversationStep) (*ai.ValidatedAnswer, error) {
	system := `You are a questionaire validator who is given answers in Dutch (quite possibly with a Frisian dialect) in the following format
		{ question: <question>, purpose: <purpose>, answer: <answer> }. The answers are transscribed from audio and might be incorrectly transcribed.
		Also they might be incomplete as the user is taking a short pause.
		I want you to try and interpret the answer in relation to the question and correct them where possible. Then I want you to 
		give an estimate if the answer is complete or not. Return this in the following json format:
		{purpose: <purpose>, answer: <answer>, <complete>: <true/false>} without any padding or fluff so I can
		directly use it in my system. It's essential that you respond in the json format I have given you.`

	question := ai.ValidatedAnswer{
		Question: step.Text,
		Purpose:  step.Purpose,
		Answer:   answer,
	}
	questionJSON, err := json.Marshal(question)
	if err != nil {
		log.Printf("Error marshaling question to JSON: %v", err)
		return nil, err
	}

	reply, err := c.AI.GetSimpleChatCompletion(system, string(questionJSON))
	if err != nil {
		log.Printf("Error getting reply from AI: %v", err)
		return nil, err
	}

	validatedAnswer := ai.ValidatedAnswer{}
	err = json.Unmarshal([]byte(reply), &validatedAnswer)
	if err != nil {
		log.Printf("Error unmarshaling validated answer: %v", err)
		return nil, err
	}

	return &validatedAnswer, nil
}

func (c *Controller) broadcastNextStep(conversationID string, state *models.ClientState, step *models.ConversationStep) (chan bool, chan error) {
	log.Println("Broadcasting next step")
	var doneChan chan bool
	var errChan chan error
	if step.AudioURL != "" {
		c.CallProvider.PlayAudioUrl(conversationID, step, state)
	} else if step.Prompt != nil {
		// TODO: implement speak from prompt
	} else {
		c.CallProvider.SpeakText(conversationID, step.Text, state)
	}

	return doneChan, errChan
}

func (c *Controller) validateAndStoreAnswer(ctx context.Context, transcript string, callID string, state *models.ClientState, rules *models.ConversationRuleSet) {
	validatedAnswer, err := c.validateAnswer(transcript, &rules.Steps[state.CurrentStep])
		if err != nil {
			log.Printf("Error validating answer, storing transcript: %v", err)
			c.storeTranscription(ctx, callID, state, rules, transcript)
		} else {
			log.Println("Validating succesful, storing validated answer")
			c.storeTranscription(ctx, callID, state, rules, validatedAnswer.Answer)
		}
}