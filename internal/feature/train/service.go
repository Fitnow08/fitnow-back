package train

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/google/uuid"
	"log/slog"
)

type TrainRepository interface {
	GetAllPublicTrains(ctx context.Context) ([]*TrainDB, error)
	GetTrainByID(ctx context.Context, id uuid.UUID) (*TrainDB, error)
	CreateTrain(ctx context.Context, req CreateTrainRequest, userID uuid.UUID) (*TrainDB, error)
	UpdateTrain(ctx context.Context, id uuid.UUID, req UpdateTrainRequest, userID uuid.UUID) (*TrainDB, error)
	DeleteTrain(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetUserTrains(ctx context.Context, userID uuid.UUID) ([]*TrainDB, error)
	AddUserTrain(ctx context.Context, userID, trainID uuid.UUID) error
	RemoveUserTrain(ctx context.Context, userID, trainID uuid.UUID) error
	GetAllExercises(ctx context.Context) ([]*ExerciseDB, error)
	CreateExercise(ctx context.Context, req CreateExerciseRequest) (*ExerciseDB, error)
}

type Service struct {
	log             *slog.Logger
	trainRepository TrainRepository
}

func NewService(log *slog.Logger, trainRepository TrainRepository) *Service {
	return &Service{log: log, trainRepository: trainRepository}
}

func (s *Service) GetAllPublicTrains(ctx context.Context) ([]*domain.Train, error) {
	rows, err := s.trainRepository.GetAllPublicTrains(ctx)
	if err != nil {
		return nil, err
	}
	trains := make([]*domain.Train, 0, len(rows))
	for _, t := range rows {
		trains = append(trains, dbToDomain(t))
	}
	return trains, nil
}

func (s *Service) GetTrainByID(ctx context.Context, id uuid.UUID) (*domain.Train, error) {
	t, err := s.trainRepository.GetTrainByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dbToDomain(t), nil
}

func (s *Service) CreateTrain(ctx context.Context, req CreateTrainRequest, userID uuid.UUID) (*domain.Train, error) {
	t, err := s.trainRepository.CreateTrain(ctx, req, userID)
	if err != nil {
		return nil, err
	}
	return dbToDomain(t), nil
}

func (s *Service) UpdateTrain(ctx context.Context, id uuid.UUID, req UpdateTrainRequest, userID uuid.UUID) (*domain.Train, error) {
	t, err := s.trainRepository.UpdateTrain(ctx, id, req, userID)
	if err != nil {
		return nil, err
	}
	return dbToDomain(t), nil
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
		trains = append(trains, dbToDomain(t))
	}
	return trains, nil
}

func (s *Service) AddUserTrain(ctx context.Context, userID, trainID uuid.UUID) error {
	return s.trainRepository.AddUserTrain(ctx, userID, trainID)
}

func (s *Service) RemoveUserTrain(ctx context.Context, userID, trainID uuid.UUID) error {
	return s.trainRepository.RemoveUserTrain(ctx, userID, trainID)
}

func (s *Service) GetAllExercises(ctx context.Context) ([]*domain.Exercise, error) {
	rows, err := s.trainRepository.GetAllExercises(ctx)
	if err != nil {
		return nil, err
	}
	exercises := make([]*domain.Exercise, 0, len(rows))
	for _, e := range rows {
		exercises = append(exercises, &domain.Exercise{
			ID:          e.ID,
			Title:       e.Title,
			Description: e.Description,
		})
	}
	return exercises, nil
}

func (s *Service) CreateExercise(ctx context.Context, req CreateExerciseRequest) (*domain.Exercise, error) {
	e, err := s.trainRepository.CreateExercise(ctx, req)
	if err != nil {
		return nil, err
	}
	return &domain.Exercise{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
	}, nil
}

func dbToDomain(t *TrainDB) *domain.Train {
	return &domain.Train{
		ID:         t.ID,
		Title:      t.Title,
		Type:       t.Type,
		Duration:   t.Duration,
		IsPublic:   t.IsPublic,
		Difficulty: t.Difficulty,
		Calories:   t.Calories,
		CreatedBy:  t.CreatedBy,
		CreatedAt:  t.CreatedAt,
	}
}
