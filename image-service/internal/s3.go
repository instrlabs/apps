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

func (s *S3Service) Upload(fileHeader *multipart.FileHeader, fileName string) error {
	file, _ := fileHeader.Open()
	defer file.Close()

	fileBytes, _ := io.ReadAll(file)

	ext := filepath.Ext(fileHeader.Filename)
	ctx := context.Background()
	_, err := s.client.PutObject(
		ctx,
		s.cfg.S3Bucket,
		fmt.Sprintf("images/%s%s", fileName, ext),
		bytes.NewReader(fileBytes),
		fileHeader.Size,
		minio.PutObjectOptions{})

	if err != nil {
		return err
	}

	return nil
}
