package db

import (
	"context"
	"goVoice/internal/config"
	"goVoice/internal/models"
	"goVoice/pkg/db/firestore"
)

type DbProvider interface {
	// Ruleset handlers
	AddRuleset(ctx context.Context, ruleset *models.ConversationRuleSet) error
	GetRuleSet(ctx context.Context, rulesetID string) (*models.ConversationRuleSet, error)
	// Conversation handlers
	GetResponses(ctx context.Context, rulesetID string, conversationID string) ([]models.ConversationStepResponse, error)
	AddConversation(ctx context.Context, rulesetID string, conversation *models.Conversation) error
	AddResponse(ctx context.Context, rulesetID string, conversationID string, response *models.ConversationStepResponse) error
	SetRecordings(ctx context.Context, rulesetID string, conversationID string, recording *models.Recording) error
	SetConversationDone(ctx context.Context, rulesetID string, conversationID string) error
	IsConversationComplete(ctx context.Context, rulesetID string, conversationID string) (*models.Recording, error)
}

func InitiateDBClient(cfg *config.Config) (DbProvider, error) {
	provider, err := firestore.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
