package app

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/config"
	constants "github.com/Sanchir01/fitnow/internal/models/contants"
	"github.com/Sanchir01/fitnow/pkg/db/connect"
	"github.com/minio/minio-go/v7"
)

type S3 struct {
	TrainBucket   *connect.Bucket
	ProgramBucket *connect.Bucket
}

func NewS3(ctx context.Context, client *minio.Client, cfg *config.Config) (*S3, error) {
	publicURL := cfg.MINIOS3.PublicBaseURL()
	trainBucket, err := connect.NewBucket(ctx, client, constants.TrainsBacket, publicURL)
	if err != nil {
		return nil, err
	}
	programBucket, err := connect.NewBucket(ctx, client, constants.ProgramsBacket, publicURL)
	if err != nil {
		return nil, err
	}
	return &S3{TrainBucket: trainBucket, ProgramBucket: programBucket}, nil
}
