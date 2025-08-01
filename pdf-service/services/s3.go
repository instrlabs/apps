package services

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

	"pdf-service/constants"
)

type S3Service struct {
	client *minio.Client
	cfg    *constants.Config
}

func NewS3Service(cfg *constants.Config) (*S3Service, error) {
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
