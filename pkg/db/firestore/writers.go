package firestore

import (
	"context"
	"goVoice/internal/models"
)

// WriteResponseToConversation writes a conversation step response to Firestore.
// It takes a context, rulesetId, conversationId, and a pointer to a ConversationStepResponse struct as input.
// It returns an error if the write operation fails.
func (f *FirestoreClient) WriteResponseToConversation(ctx *context.Context, rulesetId string, conversationId string, response *models.ConversationStepResponse) error {
	_, err := f.Client.Collection("rulesets").
		Doc(rulesetId).
		Collection("conversations").
		Doc(conversationId).
		Collection("responses").
		NewDoc().
		Create(*ctx, response)

	if err != nil {
		return err
	}
	return nil
}
