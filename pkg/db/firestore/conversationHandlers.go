package firestore

import (
	"context"
	"goVoice/internal/models"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// MUTATORS
func (f *FirestoreClient) AddConversation(ctx context.Context, rulesetID string, conversation *models.Conversation) error {
	res, err := f.Client.Collection("rulesets").
		Doc(rulesetID).
		Collection("conversations").
		Doc(conversation.ID).
		Set(ctx, conversation)

	log.Printf("Wrote to database: %v, %v", conversation.ID, res)
		
	if err != nil {
		log.Printf("Error writing conversation to firestore: %v", err)
		return err
	}
	return nil
}

func (f *FirestoreClient) AddResponse(ctx context.Context, rulesetID string, conversationID string, response *models.ConversationStepResponse) error {
	log.Printf("Adding response to firestore for %v: %v", rulesetID, response)
	docref := f.Client.Collection("rulesets").
		Doc(rulesetID).
		Collection("conversations").
		Doc(conversationID)
		
	_, err := docref.Set(ctx, map[string]interface{}{
			"responses." + response.Purpose: response.Response},
		firestore.MergeAll)
		
	if err != nil {
		log.Printf("Error writing response to firestore: %v", err)
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
		log.Printf("Error writing recording to firestore: %v", err)
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
		log.Printf("Error writing conversation done to firestore: %v", err)
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
			log.Printf("Error getting responses from firestore: %v", err)
			return nil, err
		}
		
		var response models.ConversationStepResponse
		if err := doc.DataTo(&response); err != nil {
			log.Printf("Error unmarshalling response from firestore: %v", err)
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}

// IsConversationComplete checks if a conversation is complete and returns the recording URL if it is.
func (f *FirestoreClient) IsConversationComplete (ctx context.Context, rulesetId string, conversationId string) (*models.Recording, error) {
	conversationDoc, err := f.Client.Collection("rulesets").
		Doc(rulesetId).
		Collection("conversations").
		Doc(conversationId).
		Get(ctx)
		
	if err != nil {
		log.Printf("Error getting conversation from firestore: %v", err)
		return nil, err
	}
	
	var conversation models.Conversation
	if err := conversationDoc.DataTo(&conversation); err != nil {
		log.Printf("Error unmarshalling conversation from firestore: %v", err)
		return nil, err
	}
	recordingDone := conversation.Recording != nil
	
	if conversation.ConversationDone && recordingDone {
		return conversation.Recording, nil
	}
	return nil, nil
}