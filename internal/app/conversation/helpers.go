package conversation

import (
	"context"
	"goVoice/internal/models"
)

func (c *Controller) storeTranscription(ctx context.Context, callId string, state *models.ClientState, ruleSet *models.ConversationRuleSet, transcript string) {
	c.DB.AddResponse(ctx, state.RulesetID, callId, &models.ConversationStepResponse{
		Purpose: 		ruleSet.Steps[state.CurrentStep].Purpose,
		Response: 		transcript,
	})
}
