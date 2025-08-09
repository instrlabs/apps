package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Service struct {
	client *minio.Client
	cfg    *Config
}

func NewS3Service(cfg *Config) (*S3Service, error) {
	client, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: cfg.S3UseSSL,
		Region: cfg.S3Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

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

func (s *S3Service) UploadPDF(ctx context.Context, fileHeader *multipart.FileHeader, jobID string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".pdf" {
		return "", fmt.Errorf("invalid file type: %s, expected .pdf", ext)
	}

	objectName := fmt.Sprintf("pdfs/%s%s", jobID, ext)

	_, err = s.client.PutObject(ctx, s.cfg.S3Bucket, objectName, bytes.NewReader(fileBytes), fileHeader.Size,
		minio.PutObjectOptions{ContentType: "application/pdf"})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return objectName, nil
}

// DownloadPDF downloads a PDF file from S3 to a local temporary file
func (s *S3Service) DownloadPDF(ctx context.Context, s3Path string) (string, error) {
	// Create a temporary file to store the downloaded PDF
	tempFile, err := os.CreateTemp("", "pdf-*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	tempFilePath := tempFile.Name()

	// Download the file from S3
	err = s.client.FGetObject(ctx, s.cfg.S3Bucket, s3Path, tempFilePath, minio.GetObjectOptions{})
	if err != nil {
		// Clean up the temp file if download fails
		os.Remove(tempFilePath)
		return "", fmt.Errorf("failed to download file from S3: %w", err)
	}

	return tempFilePath, nil
}

// UploadProcessedFile uploads a processed file back to S3
func (s *S3Service) UploadProcessedFile(ctx context.Context, filePath, objectName string, contentType string) (string, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = s.client.PutObject(ctx, s.cfg.S3Bucket, objectName, file, fileInfo.Size(),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", fmt.Errorf("failed to upload processed file: %w", err)
	}

	return objectName, nil
}
