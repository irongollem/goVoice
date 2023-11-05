package firestore

import (
	"context"
	"goVoice/internal/models"
)

func (f *FirestoreClient) GetRuleSet(ctx context.Context, rulesetId string) (*models.ConversationRuleSet, error) {
	docsnap, err := f.Client.Collection("rulesets").Doc(rulesetId).Get(ctx)
	if err != nil {
		return nil, err
	}
	var ruleset models.ConversationRuleSet
	if err := docsnap.DataTo(&ruleset); err != nil {
		return nil, err
	}
	return &ruleset, nil
}
