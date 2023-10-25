package firestore

import (
	"context"
	"goVoice/internal/config"

	"cloud.google.com/go/firestore"
)

type FirestoreClient struct {
	Client *firestore.Client
}

func NewClient (cfg *config.Config) (*FirestoreClient, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GCPProjectID)
	if err != nil {
		return nil, err
	}
	return &FirestoreClient{Client: client}, nil
}

