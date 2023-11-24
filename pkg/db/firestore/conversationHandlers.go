package firestore

import (
	"context"
	"goVoice/internal/models"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// MUTATORS

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

func (f *FirestoreClient) SetRecordings(ctx context.Context, rulesetId string, conversationId string, recording *models.Recording) error {
	_, err := f.Client.Collection("rulesets").
		Doc(rulesetId).
		Collection("conversations").
		Doc(conversationId).
		Update(ctx, []firestore.Update{
			{Path: "recordings", Value: recording},
		})
		
	if err != nil {
		return err
	}
	return nil
}

func (f *FirestoreClient) SetConversationDone(ctx context.Context, rulesetId string, conversationId string) error {
	_, err := f.Client.Collection("rulesets").
		Doc(rulesetId).
		Collection("conversations").
		Doc(conversationId).
		Update(ctx, []firestore.Update{
			{Path: "conversationDone", Value: true},
		})
		
	if err != nil {
		return err
	}
	return nil
}

// GETTERS

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

// IsConversationComplete checks if a conversation is complete and returns the recording URL if it is.
func (f *FirestoreClient) IsConversationComplete (ctx context.Context, rulesetId string, conversationId string) ([]models.Recording, error) {
	doc, err := f.Client.Collection("rulesets").
		Doc(rulesetId).
		Collection("conversations").
		Doc(conversationId).
		Get(ctx)
		
	if err != nil {
		log.Printf("Error getting conversation from firestore: %v", err)
		return nil, err
	}
	
	var conversation models.Conversation
	if err := doc.DataTo(&conversation); err != nil {
		log.Printf("Error unmarshalling conversation from firestore: %v", err)
		return nil, err
	}
	recordingDone := len(conversation.Recordings) > 0
	
	if conversation.ConversationDone && recordingDone {
		return conversation.Recordings, nil
	}
	return nil, nil
}