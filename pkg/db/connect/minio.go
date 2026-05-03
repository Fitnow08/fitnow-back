package connect

import (
	"context"
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
}
type Bucket struct {
	client *minio.Client
	name   string
}

func NewMinioClient(endpoint, accessKey, secretKey string, useSSL bool) (*minio.Client, error) {
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
}

func NewBucket(ctx context.Context, client *minio.Client, name string) (*Bucket, error) {
	exists, err := client.BucketExists(ctx, name)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, name, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}
	return &Bucket{client: client, name: name}, nil
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
