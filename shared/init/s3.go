package initx

import (
	"bytes"
	"context"
	"io"

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

func NewS3(cfg *S3Config) (*S3, error) {
	client, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: cfg.S3UseSSL,
		Region: cfg.S3Region,
	})
	if err != nil {
		return nil, err
	}
	_ = client.MakeBucket(context.Background(), cfg.S3Bucket, minio.MakeBucketOptions{Region: cfg.S3Region})
	return &S3{Client: client, Cfg: cfg}, nil
}

func (s *S3) Put(objectName string, data []byte, contentType string) error {
	_, err := s.Client.PutObject(context.Background(), s.Cfg.S3Bucket, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (s *S3) Get(objectName string) ([]byte, error) {
	r, err := s.Client.GetObject(context.Background(), s.Cfg.S3Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer r.Close()
	b, err := io.ReadAll(r)
	return b, err
}
