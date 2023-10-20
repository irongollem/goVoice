package storage

import (
	"goVoice/internal/config"

	"io"

	"github.com/gin-gonic/gin"
)

type StorageProvider interface {
	UploadFile(ctx *gin.Context, filename string, file io.Reader) error
	DownloadFile(ctx *gin.Context, filename string) (io.ReadCloser, error)
	DeleteFile(ctx *gin.Context, filename string) error
}

func NewStorageHandler(cfg *config.Config) (StorageProvider, error) {
	googleProvider, _ := NewGoogleStorageHandler(cfg)

	return googleProvider, nil
}
