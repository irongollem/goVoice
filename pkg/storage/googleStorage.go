package storage

import (
	"context"
	"fmt"
	"goVoice/internal/config"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type GoogleStorageHandler struct {
	client *storage.Client
	bucket string
}

func NewGoogleStorageHandler(cfg *config.Config) (*GoogleStorageHandler, error) {
	ctx := context.Background()
	bucket := "govoice-recordings"

	credentials, err := os.ReadFile(cfg.GCPCredentialsFile)
	if err != nil {
		return nil, err
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credentials))
	if err != nil {
		return nil, err
	}
	
	return &GoogleStorageHandler{
		client: client,
		bucket: bucket,
	}, nil
}

func (h *GoogleStorageHandler) GetRecording(ctx context.Context, rulesetId string, callId string) (io.ReadCloser, error) {
	h.DownloadFile(ctx, fmt.Sprintf("recording-%s-%s", rulesetId, callId))
	return nil, nil
}

func (h *GoogleStorageHandler) UploadFile(ctx context.Context, filename string, file io.Reader) error {
	sw := h.client.Bucket(h.bucket).Object(filename).NewWriter(ctx)
	if _, err := io.Copy(sw, file); err != nil {
		log.Printf("Error uploading file to Google Cloud Storage: %v", err)
		return err
	}
	if err := sw.Close(); err != nil {
		return err
	}
	return nil
}

func (h *GoogleStorageHandler) DownloadFile(ctx context.Context, filename string) (io.ReadCloser, error) {
	reader, err := h.client.Bucket(h.bucket).Object(filename).NewReader(ctx)
	if err != nil {
		log.Printf("Error downloading file from Google Cloud Storage: %v", err)
		return nil, err
	}
	return reader, nil
}

func (h *GoogleStorageHandler) DeleteFile(c *gin.Context, filename string) error {
	if err := h.client.Bucket(h.bucket).Object(filename).Delete(c); err != nil {
		log.Printf("Error deleting file from Google Cloud Storage: %v", err)
		return err
	}
	return nil
}

func (h *GoogleStorageHandler) Close() error {
	if err := h.client.Close(); err != nil {
		log.Printf("Failed to close Google Cloud Storage client: %v", err)
		return err
	}
	return nil
}

func (h *GoogleStorageHandler) MoveRenameFile(ctx context.Context, srcName string, dstName string) error {
	src := h.client.Bucket(h.bucket).Object(srcName)
	dst := h.client.Bucket(h.bucket).Object(dstName)
	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		log.Printf("Error renaming file in Google Cloud Storage: %v", err)
		return err
	}
	if err := src.Delete(ctx); err != nil {
		log.Printf("Error deleting file from Google Cloud Storage: %v", err)
		return err
	}
	return nil
}
