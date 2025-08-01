package services

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"labs-worker/constants"
)

// S3Service handles interactions with S3-compatible storage
type S3Service struct {
	client *minio.Client
	cfg    *constants.Config
}

// NewS3Service creates a new S3Service
func NewS3Service(cfg *constants.Config) (*S3Service, error) {
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

// DownloadPDF downloads a PDF file from S3 to a local file
func (s *S3Service) DownloadPDF(ctx context.Context, objectName string) (string, error) {
	// Create a temporary file to store the PDF
	tempFile, err := os.CreateTemp("", "pdf-*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	// Download the object
	err = s.client.FGetObject(ctx, s.cfg.S3Bucket, objectName, tempFile.Name(), minio.GetObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}

	return tempFile.Name(), nil
}

// UploadJPG uploads a JPG file to S3
func (s *S3Service) UploadJPG(ctx context.Context, filePath string, jobID string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Create the object name (path in the bucket)
	objectName := fmt.Sprintf("jpgs/%s.jpg", jobID)

	// Upload the file
	_, err = s.client.PutObject(ctx, s.cfg.S3Bucket, objectName, file, fileInfo.Size(),
		minio.PutObjectOptions{ContentType: "image/jpeg"})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return objectName, nil
}

// UploadFromBytes uploads a file from bytes to S3
func (s *S3Service) UploadFromBytes(ctx context.Context, data []byte, objectName string, contentType string) error {
	// Upload the file
	_, err := s.client.PutObject(ctx, s.cfg.S3Bucket, objectName, bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}
