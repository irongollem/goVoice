package db

import (
	"context"
	"goVoice/internal/config"
	"goVoice/internal/models"
	"goVoice/pkg/db/firestore"
)

type DbProvider interface {
	GetRuleSet(ctx context.Context, rulesetId string) (*models.ConversationRuleSet, error)
	GetResponses(ctx context.Context, rulesetId string, conversationId string) ([]models.ConversationStepResponse, error)

	AddConversation(ctx context.Context, rulesetId string, conversation *models.Conversation) error
	AddResponse(ctx context.Context, rulesetId string, conversationId string, response *models.ConversationStepResponse) error
}

func InitiateDBClient(cfg *config.Config) (DbProvider, error) {
	provider, err := firestore.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
