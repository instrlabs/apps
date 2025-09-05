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

// UploadBytes uploads in-memory data to S3 under the provided objectName and contentType.
func (s *S3Service) UploadBytes(data []byte, objectName, contentType string) error {
	ctx := context.Background()
	_, err := s.client.PutObject(
		ctx,
		s.cfg.S3Bucket,
		objectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	return err
}

func (s *S3Service) DownloadAllForInstruction(instructionID string) ([][]byte, []string, error) {
	ctx := context.Background()
	prefix := fmt.Sprintf("images/%s-", instructionID)
	objCh := s.client.ListObjects(ctx, s.cfg.S3Bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true})

	var datas [][]byte
	var keys []string
	for obj := range objCh {
		if obj.Err != nil {
			return nil, nil, obj.Err
		}
		reader, err := s.client.GetObject(ctx, s.cfg.S3Bucket, obj.Key, minio.GetObjectOptions{})
		if err != nil {
			return nil, nil, err
		}
		b, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			return nil, nil, err
		}
		datas = append(datas, b)
		keys = append(keys, obj.Key)
	}
	if len(datas) == 0 {
		return nil, nil, fmt.Errorf("no request file found for instruction %s", instructionID)
	}
	return datas, keys, nil
}
