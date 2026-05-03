package app

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/config"
	constants "github.com/Sanchir01/fitnow/internal/models/contants"
	"github.com/Sanchir01/fitnow/pkg/db/connect"
	"github.com/minio/minio-go/v7"
)

type S3 struct {
	TrainBucket *connect.Bucket
}

func NewS3(ctx context.Context, client *minio.Client, cfg *config.Config) (*S3, error) {
	trainBucket, err := connect.NewBucket(ctx, client, constants.TrainsBacket)
	if err != nil {
		return nil, err
	}
	return &S3{TrainBucket: trainBucket}, nil
}
