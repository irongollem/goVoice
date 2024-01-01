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
	for purpose, answer := range responses {
		sb.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n", rulesetID, rulesetTitle, callID, purpose, answer))
	}
	return sb.String()
}

func (c *Controller) validateAnswer(answer string, step *models.ConversationStep) (*ai.ValidatedAnswer, error) {
	system := `You are a questionaire validator who is given answers in Dutch in the following format
		{ question: <question>, purpose: <purpose>, answer: <answer> }. The answers are transscribed from audio and might be incorrectly transcribed.
		Also they might be incomplete as the user is taking a short pause.
		I want you to try and interpret the answer and correct them where possible. Then I want you to 
		answer give an estimate if the answer is complete or not. Return this in the following format:
		{purpose: <purpose>, answer: <answer>, <complete>: <true/false>} without any padding or fluff so I can
		directly use it in my system.`

	question := ai.ValidatedAnswer{
		Question: step.Text,
		Purpose: step.Purpose,
		Answer: answer,
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