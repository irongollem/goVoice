package firestore

import (
	"context"
	"goVoice/internal/models"

	"cloud.google.com/go/firestore"
)

func (f *FirestoreClient) AddConversation(ctx context.Context, rulesetId string, conversation *models.Conversation) error {
	_, err := f.Client.Collection("rulesets").
		Doc(rulesetId).
		Collection("conversations").
		Doc(conversation.ID).
		Set(ctx, conversation)
		
	if err != nil {
		return err
	}
	return nil
}

func (f *FirestoreClient) AddResponse(ctx context.Context, rulesetId string, conversationId string, response *models.ConversationStepResponse) error {
	_, err := f.Client.Collection("rulesets").
		Doc(rulesetId).
		Collection("conversations").
		Doc(conversationId).
		Update(ctx, []firestore.Update{
			{Path: "responses." + response.Purpose, Value: response.Response},
		})
		
	if err != nil {
		return err
	}
	return nil
}