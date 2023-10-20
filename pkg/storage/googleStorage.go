package storage

import (
	"context"
	"goVoice/internal/config"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
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

func (h *GoogleStorageHandler) UploadFile(c *gin.Context, filename string, file io.Reader) error {
	ctx := appengine.NewContext(c.Request)
	sw := h.client.Bucket(h.bucket).Object(filename).NewWriter(ctx)
	if _, err := io.Copy(sw, file); err != nil {
		return err
	}
	if err := sw.Close(); err != nil {
		return err
	}
	return nil
}

func (h *GoogleStorageHandler) DownloadFile(c *gin.Context, filename string) (io.ReadCloser, error) {
	rc, err := h.client.Bucket(h.bucket).Object(filename).NewReader(c)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

func (h *GoogleStorageHandler) DeleteFile(c *gin.Context, filename string) error {
	if err := h.client.Bucket(h.bucket).Object(filename).Delete(c); err != nil {
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
