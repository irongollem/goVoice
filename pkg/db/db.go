package db

import (
	"context"
	"goVoice/internal/config"
)



type DbProvider interface {
	GetDoc(ctx context.Context, rulesetId string) (*Ruleset, error)
}

func NewDbHandler(cfg *config.Config) (*DbProvider, error) {
	newFirestoreHandler, err := NewFirestoreHandler(cfg)
	return newFirestoreHandler, err
}