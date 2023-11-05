package db

import (
	"context"
	"goVoice/internal/config"
	"goVoice/internal/models"
	"goVoice/pkg/db/firestore"
)

type DbProvider interface {
	// Ruleset handlers
	GetRuleSet(ctx context.Context, rulesetId string) (*models.ConversationRuleSet, error)
	// Conversation handlers
	GetResponses(ctx context.Context, rulesetId string, conversationId string) ([]models.ConversationStepResponse, error)
	AddConversation(ctx context.Context, rulesetId string, conversation *models.Conversation) error
	AddResponse(ctx context.Context, rulesetId string, conversationId string, response *models.ConversationStepResponse) error
	SetRecordings(ctx context.Context, rulesetId string, conversationId string, recordings []models.Recording) error
	SetConversationDone(ctx context.Context, rulesetId string, conversationId string) error
	IsConversationComplete(ctx context.Context, rulesetId string, conversationId string) ([]models.Recording, error)
}

func InitiateDBClient(cfg *config.Config) (DbProvider, error) {
	provider, err := firestore.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
