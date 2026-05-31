package connect

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"net/url"
	"time"
)

type MiniS3Interface interface {
	Upload(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	Delete(ctx context.Context, key string) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	PresignedGetURL(ctx context.Context, key string, ttl time.Duration) (*url.URL, error)
	PublicURL(key string) string
}
type Bucket struct {
	client    *minio.Client
	name      string
	publicURL string
}

func NewMinioClient(endpoint, accessKey, secretKey string, useSSL bool) (*minio.Client, error) {
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
}

func NewBucket(ctx context.Context, client *minio.Client, name, publicURL string) (*Bucket, error) {
	exists, err := client.BucketExists(ctx, name)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, name, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}
	// public-read: объекты доступны по прямой ссылке без подписи (для отдачи картинок).
	if err := client.SetBucketPolicy(ctx, name, publicReadPolicy(name)); err != nil {
		return nil, err
	}
	return &Bucket{client: client, name: name, publicURL: publicURL}, nil
}

func publicReadPolicy(bucket string) string {
	return fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucket)
}

// PublicURL собирает прямую ссылку на объект: <publicURL>/<bucket>/<key>.
// Пустой key даёт пустую строку.
func (b *Bucket) PublicURL(key string) string {
	if key == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s", b.publicURL, b.name, key)
}

func (b *Bucket) Upload(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	_, err := b.client.PutObject(ctx, b.name, key, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (b *Bucket) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	return b.client.GetObject(ctx, b.name, key, minio.GetObjectOptions{})
}

func (b *Bucket) Delete(ctx context.Context, key string) error {
	return b.client.RemoveObject(ctx, b.name, key, minio.RemoveObjectOptions{})
}

func (b *Bucket) PresignedGetURL(ctx context.Context, key string, ttl time.Duration) (*url.URL, error) {
	return b.client.PresignedGetObject(ctx, b.name, key, ttl, nil)
}
