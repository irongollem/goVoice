package storage

import (
	"context"
	"goVoice/internal/config"

	"io"
)

type StorageProvider interface {
	// UploadFile(ctx *gin.Context, filename string, file io.Reader) error
	// DownloadFile(ctx *gin.Context, filename string) (io.ReadCloser, error)
	// DeleteFile(ctx *gin.Context, filename string) error
	GetRecording(ctx context.Context, rulesetID string, callID string) (io.ReadCloser, error)
}

func NewStorageHandler(cfg *config.Config) (StorageProvider, error) {
	provider, _ := NewGoogleStorageHandler(cfg)

	return provider, nil
}
