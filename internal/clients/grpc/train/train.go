package traingrpc

import (
	"context"
	"crypto/tls"
	"fmt"
	trainv1 "github.com/Fitnow08/fitnow-proto/pkg/gen/go/v1/train"
	"github.com/Sanchir01/fitnow/internal/config"
	"github.com/Sanchir01/fitnow/pkg/logger"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type TrainClient struct {
	trainv1.TrainServiceClient
	log *slog.Logger
}

func NewTrainClient(log *slog.Logger, cfg config.Client) (*TrainClient, error) {
	const op = "TrainClient.NewTrainClient"
	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Aborted, codes.DeadlineExceeded, codes.NotFound),
		grpcretry.WithMax(uint(cfg.Retries)),
		grpcretry.WithPerRetryTimeout(cfg.Timeout),
	}

	creds := insecure.NewCredentials()
	if !cfg.Insecure {
		creds = credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
	}

	cc, err := grpc.NewClient(
		cfg.Address,
		grpc.WithTransportCredentials(creds),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
			grpclog.UnaryClientInterceptor(logger.InterceptorLogger(log), logOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &TrainClient{
		TrainServiceClient: trainv1.NewTrainServiceClient(cc),
		log:                log,
	}, nil
}

func (c *TrainClient) GetTrainAndExercises(ctx context.Context, req *trainv1.GetTrainAndExercisesRequest) (*trainv1.GetTrainAndExercisesResponse, error) {
	const op = "TrainClient.GetTrainAndExercises"
	resp, err := c.TrainServiceClient.GetTrainAndExercises(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *TrainClient) AddTrainExercises(ctx context.Context, req *trainv1.AddTrainExercisesRequest) (*trainv1.AddTrainExercisesResponse, error) {
	const op = "TrainClient.AddTrainExercises"
	resp, err := c.TrainServiceClient.AddTrainExercises(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *TrainClient) GetAllAdminTrains(ctx context.Context, req *trainv1.GetAllAdminTrainsRequest) (*trainv1.GetAllAdminTrainsResponse, error) {
	const op = "TrainClient.GetAllAdminTrains"
	resp, err := c.TrainServiceClient.GetAllAdminTrains(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *TrainClient) GetAllTrains(ctx context.Context, req *trainv1.GetAllTrainsRequest) (*trainv1.GetAllTrainsResponse, error) {
	const op = "TrainClient.GetAllTrains"
	resp, err := c.TrainServiceClient.GetAllTrains(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *TrainClient) GetTrainById(ctx context.Context, req *trainv1.GetTrainByIdRequest) (*trainv1.GetTrainByIdResponse, error) {
	const op = "TrainClient.GetTrainById"
	resp, err := c.TrainServiceClient.GetTrainById(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *TrainClient) CreateTrain(ctx context.Context, req *trainv1.CreateTrainRequest) (*trainv1.CreateTrainResponse, error) {
	const op = "TrainClient.CreateTrain"
	resp, err := c.TrainServiceClient.CreateTrain(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *TrainClient) UpdateTrain(ctx context.Context, req *trainv1.UpdateTrainRequest) (*trainv1.UpdateTrainResponse, error) {
	const op = "TrainClient.UpdateTrain"
	resp, err := c.TrainServiceClient.UpdateTrain(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *TrainClient) DeleteTrain(ctx context.Context, req *trainv1.DeleteTrainRequest) (*trainv1.DeleteTrainResponse, error) {
	const op = "TrainClient.DeleteTrain"
	resp, err := c.TrainServiceClient.DeleteTrain(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *TrainClient) AddTrainImage(ctx context.Context, req *trainv1.AddTrainImageRequest) (*trainv1.AddTrainImageResponse, error) {
	const op = "TrainClient.AddTrainImage"
	resp, err := c.TrainServiceClient.AddTrainImage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}
