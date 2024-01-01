package firestore

import (
	"context"
	"goVoice/internal/models"
	"log"

	"cloud.google.com/go/firestore"
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

func (f *FirestoreClient) DeleteConversation(ctx context.Context, rulesetID string, conversationID string) error {
	_, err := f.Client.Collection("rulesets").
		Doc(rulesetID).
		Collection("conversations").
		Doc(conversationID).
		Delete(ctx)
		
	if err != nil {
		log.Printf("Error deleting conversation from firestore: %v", err)
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
		
	_, err := docref.Update(ctx, []firestore.Update{
		{
			Path:  "responses." + response.Purpose,
			Value: response.Response,
		},
	})
		
	if err != nil {
		log.Printf("Error writing response to firestore: %v", err)
		return err
	}
	return nil
}

func (f *FirestoreClient) SetRecording(ctx context.Context, rulesetId string, conversationId string, recording *models.Recording) error {
	_, err := f.Client.Collection("rulesets").
		Doc(rulesetId).
		Collection("conversations").
		Doc(conversationId).
		Update(ctx, []firestore.Update{
			{Path: "recordings", Value:  firestore.ArrayUnion((recording))},
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

func (f *FirestoreClient) GetConversation(ctx context.Context, rulesetId string, conversationId string) (*models.Conversation, error) {
	conversation := f.Client.Collection("rulesets").
		Doc(rulesetId).
		Collection("conversations").
		Doc(conversationId)

	docsnap, err := conversation.Get(ctx)
	if err != nil {
		log.Printf("Error getting conversation from firestore: %v", err)
		return nil, err
	}

	var c models.Conversation
	err = docsnap.DataTo(&c)
	if err != nil {
		log.Printf("Error unmarshalling responses from firestore: %v", err)
		return nil, err
	}

	return &c, nil
}

func (f *FirestoreClient) GetRecordings(ctx context.Context, rulesetID string, conversationID string) ([]models.Recording, error) {
	conversationDoc, err := f.Client.Collection("rulesets").
		Doc(rulesetID).
		Collection("conversations").
		Doc(conversationID).
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
	
	return conversation.Recordings, nil
}