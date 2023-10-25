package db

import (
	"context"
	"goVoice/internal/config"
	"goVoice/internal/models"
	"goVoice/pkg/db/firestore"
)

type DbProvider interface {
	GetRuleSet(ctx *context.Context, rulesetId string) (*models.ConversationRuleSet, error)
}

func InitiateDBClient(cfg *config.Config) (DbProvider, error) {
	provider, err := firestore.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
