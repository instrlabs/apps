package initx

import (
	"bytes"
	"context"
	"io"

	"github.com/gofiber/fiber/v2/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Config struct {
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3UseSSL    bool
	S3Region    string
	S3Bucket    string
}

type S3 struct {
	Client *minio.Client
	Cfg    *S3Config
}

func NewS3(cfg *S3Config) *S3 {
	client, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: cfg.S3UseSSL,
		Region: cfg.S3Region,
	})
	if err != nil {
		log.Errorf("Failed to create S3 client: %v", err)
		return nil
	}
	_ = client.MakeBucket(context.Background(), cfg.S3Bucket, minio.MakeBucketOptions{Region: cfg.S3Region})
	return &S3{Client: client, Cfg: cfg}
}

func (s *S3) Put(objectName string, data []byte) error {
	_, err := s.Client.PutObject(context.Background(), s.Cfg.S3Bucket, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		log.Errorf("Failed to put object %s to S3: %v", objectName, err)
		return err
	}
	return nil
}

func (s *S3) Get(objectName string) []byte {
	r, err := s.Client.GetObject(context.Background(), s.Cfg.S3Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Errorf("Failed to get object %s from S3: %v", objectName, err)
		return nil
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		log.Errorf("Failed to read object %s from S3: %v", objectName, err)
		return nil
	}
	return b
}
