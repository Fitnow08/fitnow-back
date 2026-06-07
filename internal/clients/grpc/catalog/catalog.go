package catalog

import (
	"context"
	"crypto/tls"
	"fmt"
	programv1 "github.com/Fitnow08/fitnow-proto/pkg/gen/go/v1/program"
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

type ProgramClient struct {
	programv1.ProgramServiceClient
	log *slog.Logger
}

func NewProgramClient(log *slog.Logger, cfg config.Client) (*ProgramClient, error) {
	const op = "ProgramClient.NewProgramClient"
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

	return &ProgramClient{
		ProgramServiceClient: programv1.NewProgramServiceClient(cc),
		log:                  log,
	}, nil
}

func (c *ProgramClient) GetAllPrograms(ctx context.Context, categoryID, search string) (*programv1.GetAllProgramsResponse, error) {
	const op = "ProgramClient.GetAllPrograms"
	resp, err := c.ProgramServiceClient.GetAllPrograms(ctx, &programv1.GetAllProgramsRequest{
		CategoryId: categoryID,
		Search:     search,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *ProgramClient) CreateProgram(ctx context.Context, req *programv1.CreateProgramRequest) (*programv1.CreateProgramResponse, error) {
	const op = "ProgramClient.CreateProgram"
	resp, err := c.ProgramServiceClient.CreateProgram(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *ProgramClient) AddProgramImage(ctx context.Context, req *programv1.AddProgramImageRequest) (*programv1.AddProgramImageResponse, error) {
	const op = "ProgramClient.AddProgramImage"
	resp, err := c.ProgramServiceClient.AddProgramImage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *ProgramClient) GetProgramsAndTrains(ctx context.Context, programID string) (*programv1.GetProgramsAndTrainsResponse, error) {
	const op = "ProgramClient.GetProgramsAndTrains"
	resp, err := c.ProgramServiceClient.GetProgramsAndTrains(ctx, &programv1.GetProgramsAndTrainsRequest{
		ProgramId: programID,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *ProgramClient) UploadAllProgramTrains(ctx context.Context, programID string, trains []*programv1.ProgramTrainInput) (*programv1.UploadAllProgramTrainsResponse, error) {
	const op = "ProgramClient.UploadAllProgramTrains"
	resp, err := c.ProgramServiceClient.UploadAllProgramTrains(ctx, &programv1.UploadAllProgramTrainsRequest{
		ProgramId: programID,
		Trains:    trains,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}
