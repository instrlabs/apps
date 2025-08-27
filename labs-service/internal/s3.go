package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Service handles interactions with S3-compatible storage
type S3Service struct {
	client *minio.Client
	cfg    *Config
}

// NewS3Service creates a new S3Service
func NewS3Service(cfg *Config) (*S3Service, error) {
	// Initialize minio client
	client, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: cfg.S3UseSSL,
		Region: cfg.S3Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	// Check if bucket exists, create if it doesn't
	exists, err := client.BucketExists(context.Background(), cfg.S3Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		err = client.MakeBucket(context.Background(), cfg.S3Bucket, minio.MakeBucketOptions{
			Region: cfg.S3Region,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Created bucket %s\n", cfg.S3Bucket)
	}

	return &S3Service{
		client: client,
		cfg:    cfg,
	}, nil
}

// UploadPDF uploads a PDF file to S3
func (s *S3Service) UploadPDF(ctx context.Context, fileHeader *multipart.FileHeader, jobID string) (string, error) {
	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read the file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Get the file extension
	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".pdf" {
		return "", fmt.Errorf("invalid file type: %s, expected .pdf", ext)
	}

	// Create the object name (path in the bucket)
	objectName := fmt.Sprintf("pdfs/%s%s", jobID, ext)

	// Upload the file
	_, err = s.client.PutObject(ctx, s.cfg.S3Bucket, objectName, bytes.NewReader(fileBytes), fileHeader.Size,
		minio.PutObjectOptions{ContentType: "application/pdf"})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return objectName, nil
}
