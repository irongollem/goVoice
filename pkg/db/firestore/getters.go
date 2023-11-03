package firestore

import (
	"context"
	"goVoice/internal/models"

	"google.golang.org/api/iterator"
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

func (f *FirestoreClient) GetResponses (ctx context.Context, rulesetId string, conversationId string) ([]models.ConversationStepResponse, error) {
	iter := f.Client.Collection("rulesets").
	Doc(rulesetId).
	Collection("conversations").
	Doc(conversationId).
	Collection("responses").
	Documents(ctx)
	
	var responses []models.ConversationStepResponse
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		
		var response models.ConversationStepResponse
		if err := doc.DataTo(&response); err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}
