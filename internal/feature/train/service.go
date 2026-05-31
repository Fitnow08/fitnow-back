package train

import (
	"context"
	"fmt"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/Sanchir01/fitnow/pkg/db/connect"
	"github.com/google/uuid"
	"io"
	"log/slog"
)

type TrainRepository interface {
	GetAllPublicTrains(ctx context.Context, param AllTrainsParams) ([]*TrainDB, error)
	GetTrainByID(ctx context.Context, id uuid.UUID) (*TrainDB, error)
	CreateTrain(ctx context.Context, req CreateTrainRequest, userID uuid.UUID) (*TrainDB, error)
	UpdateTrain(ctx context.Context, id uuid.UUID, req UpdateTrainRequest, userID uuid.UUID) (*TrainDB, error)
	DeleteTrain(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetUserTrains(ctx context.Context, userID uuid.UUID) ([]*TrainDB, error)
	AddUserTrain(ctx context.Context, userID, trainID uuid.UUID) error
	RemoveUserTrain(ctx context.Context, userID, trainID uuid.UUID) error
	UpdateTrainImageUrl(ctx context.Context, trainID uuid.UUID, url string) error
}

type Service struct {
	log             *slog.Logger
	trainRepository TrainRepository
	s3              connect.MiniS3Interface
}

func NewService(log *slog.Logger, s3 connect.MiniS3Interface, trainRepository TrainRepository) *Service {

	return &Service{log: log, trainRepository: trainRepository, s3: s3}
}

func (s *Service) GetAllPublicTrains(ctx context.Context, param AllTrainsParams) ([]*domain.Train, error) {
	rows, err := s.trainRepository.GetAllPublicTrains(ctx, param)
	if err != nil {
		return nil, err
	}
	trains := make([]*domain.Train, 0, len(rows))
	for _, t := range rows {
		trains = append(trains, s.dbToDomain(t))
	}
	return trains, nil
}

func (s *Service) GetTrainByID(ctx context.Context, id uuid.UUID) (*domain.Train, error) {
	t, err := s.trainRepository.GetTrainByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.dbToDomain(t), nil
}

func (s *Service) CreateTrain(ctx context.Context, req CreateTrainRequest, userID uuid.UUID) (*domain.Train, error) {
	t, err := s.trainRepository.CreateTrain(ctx, req, userID)
	if err != nil {
		return nil, err
	}
	return s.dbToDomain(t), nil
}

func (s *Service) UpdateTrain(ctx context.Context, id uuid.UUID, req UpdateTrainRequest, userID uuid.UUID) (*domain.Train, error) {
	t, err := s.trainRepository.UpdateTrain(ctx, id, req, userID)
	if err != nil {
		return nil, err
	}
	return s.dbToDomain(t), nil
}

func (s *Service) DeleteTrain(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.trainRepository.DeleteTrain(ctx, id, userID)
}

func (s *Service) GetUserTrains(ctx context.Context, userID uuid.UUID) ([]*domain.Train, error) {
	rows, err := s.trainRepository.GetUserTrains(ctx, userID)
	if err != nil {
		return nil, err
	}
	trains := make([]*domain.Train, 0, len(rows))
	for _, t := range rows {
		trains = append(trains, s.dbToDomain(t))
	}
	return trains, nil
}

func (s *Service) AddUserTrain(ctx context.Context, userID, trainID uuid.UUID) error {
	return s.trainRepository.AddUserTrain(ctx, userID, trainID)
}

func (s *Service) RemoveUserTrain(ctx context.Context, userID, trainID uuid.UUID) error {
	return s.trainRepository.RemoveUserTrain(ctx, userID, trainID)
}

func (s *Service) UploadTrainImage(ctx context.Context, trainID uuid.UUID, ext, contentType string, size int64, r io.Reader) error {
	key := fmt.Sprintf("trains/%s%s", trainID, ext) // например trains/<uuid>.webp
	if err := s.s3.Upload(ctx, key, r, size, contentType); err != nil {
		return err
	}
	if err := s.trainRepository.UpdateTrainImageUrl(ctx, trainID, key); err != nil {
		return err
	}
	return nil
}

func (s *Service) dbToDomain(t *TrainDB) *domain.Train {
	return &domain.Train{
		ID:         t.ID,
		Title:      t.Title,
		Type:       t.Type,
		Duration:   t.Duration,
		IsPublic:   t.IsPublic,
		Difficulty: t.Difficulty,
		CategoryId: t.CategoryId,
		Calories:   t.Calories,
		ImageURL:   s.s3.PublicURL(t.ImagePath),
		CreatedBy:  t.CreatedBy,
		CreatedAt:  t.CreatedAt,
	}
}
