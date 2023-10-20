package db

import (
	"context"
	"goVoice/internal/config"

	"cloud.google.com/go/firestore"
)

type FirestoreHandler struct {
	Client *firestore.Client
}

func NewFirestoreHandler(cfg *config.Config) (*DbProvider, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "goVoice")
	if err != nil {
		return nil, err
	}
	return &FirestoreHandler{Client: client}, nil
}

func (h *FirestoreHandler) GetDoc(ctx context.Context, rulesetId string) (*Ruleset, error) {
	doc, err := h.Client.Collection("rulesets").Doc(rulesetId).Get(ctx)
	if err != nil {
		return nil, err
	}
	var ruleset Ruleset
	doc.DataTo(&ruleset)
	return &ruleset, nil
}