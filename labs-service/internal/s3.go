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

type S3Service struct {
	client *minio.Client
	cfg    *Config
}

func NewS3Service(cfg *Config) *S3Service {
	client, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: cfg.S3UseSSL,
		Region: cfg.S3Region,
	})
	if err != nil {
		log.Printf("failed to create S3 client: %v", err)
		return &S3Service{client: nil, cfg: cfg}
	}

	exists, err := client.BucketExists(context.Background(), cfg.S3Bucket)
	if err != nil {
		log.Printf("failed to check if bucket exists: %v", err)
		return &S3Service{client: client, cfg: cfg}
	}

	if !exists {
		err = client.MakeBucket(context.Background(), cfg.S3Bucket, minio.MakeBucketOptions{
			Region: cfg.S3Region,
		})
		if err != nil {
			log.Printf("failed to create bucket: %v", err)
			return &S3Service{client: client, cfg: cfg}
		}
		log.Printf("Created bucket %s\n", cfg.S3Bucket)
	}

	return &S3Service{
		client: client,
		cfg:    cfg,
	}
}

func (s *S3Service) UploadPDF(ctx context.Context, fileHeader *multipart.FileHeader, jobID string) (string, error) {
	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("failed to open file: %v", err)
		return "", nil
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("failed to read file: %v", err)
		return "", nil
	}

	// Get the file extension
	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".pdf" {
		log.Printf("invalid file type: %s, expected .pdf", ext)
		return "", nil
	}

	objectName := fmt.Sprintf("pdfs/%s%s", jobID, ext)

	if s.client == nil {
		log.Printf("S3 client is not initialized; skipping upload for %s", objectName)
		return "", nil
	}

	_, err = s.client.PutObject(ctx, s.cfg.S3Bucket, objectName, bytes.NewReader(fileBytes), fileHeader.Size,
		minio.PutObjectOptions{ContentType: "application/pdf"})
	if err != nil {
		log.Printf("failed to upload file: %v", err)
		return "", nil
	}

	return objectName, nil
}
