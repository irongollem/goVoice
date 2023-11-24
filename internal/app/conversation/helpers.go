package conversation

import (
	"context"
	"fmt"
	"goVoice/internal/models"
	"log"
	"strings"
)

func (c *Controller) storeTranscription(ctx context.Context, callId string, state *models.ClientState, ruleSet *models.ConversationRuleSet, transcript string) {
	log.Printf("Storing transcription: %s", transcript)
	c.DB.AddResponse(ctx, state.RulesetID, callId, &models.ConversationStepResponse{
		Purpose: 		ruleSet.Steps[state.CurrentStep].Purpose,
		Response: 		transcript,
	})
}

func (c *Controller) retrieveResponsesAsTable(ctx context.Context, ruleset *models.ConversationRuleSet, callId string) (*string, error) {
		// Get the stored responses from the database
		responses, err := c.DB.GetResponses(ctx, ruleset.ID, callId)
		if err != nil {
			log.Printf("Error getting responses from database: %v", err)
			return nil, err
		}
		body := formatEmailBody(responses, ruleset.ID, ruleset.Title, callId)
	return &body, nil

}

func formatEmailBody(responses []models.ConversationStepResponse, rulesetID string, rulesetTitle string, callID string) string {
	var sb strings.Builder
	for _, response := range responses {
		sb.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n", rulesetID, rulesetTitle, callID, response.Purpose, response.Response))
	}
	return sb.String()
}
