package firestore

import (
	"context"
	"goVoice/internal/models"
	"log"
)

func (f *FirestoreClient) AddRuleset(ctx context.Context, ruleset *models.ConversationRuleSet) error {
	_, err := f.Client.Collection("rulesets").
		Doc(ruleset.ID).
		Set(ctx, ruleset)

	log.Printf("Wrote to database: %v", ruleset.Title)

	if err != nil {
		log.Printf("Error writing ruleset to firestore: %v", err)
		return err
	}
	return nil
}

func (f *FirestoreClient) GetRuleSet(ctx context.Context, rulesetID string) (*models.ConversationRuleSet, error) {
	docsnap, err := f.Client.Collection("rulesets").
		Doc(rulesetID).
		Get(ctx)

	if err != nil {
		return nil, err
	}
	var ruleset models.ConversationRuleSet
	if err := docsnap.DataTo(&ruleset); err != nil {
		return nil, err
	}
	return &ruleset, nil
}
