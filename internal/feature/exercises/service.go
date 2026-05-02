package exercises

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"log/slog"
)

type ExercisesRepository interface {
	GetAllExercises(ctx context.Context) ([]ExerciseDB, error)
	CreateExercise(ctx context.Context, title, desc string) (*ExerciseDB, error)
}

type Service struct {
	log  *slog.Logger
	repo ExercisesRepository
}

func NewService(log *slog.Logger, repo ExercisesRepository) *Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) GetAllExercises(ctx context.Context) ([]domain.Exercise, error) {
	const op = "Exercise.Service.GetAllExercises"
	log := s.log.With("op", op)
	data, err := s.repo.GetAllExercises(ctx)
	if err != nil {
		log.Error(op, "failed to get all exercises")
		return nil, err
	}
	var exercises = make([]domain.Exercise, len(data))
	for _, info := range data {
		exercises = append(exercises, domain.Exercise{
			ID:          info.ID,
			Title:       info.Title,
			Description: info.Description,
		})
	}
	return exercises, nil
}

func (s *Service) CreateExercise(ctx context.Context, title, desc string) (*domain.Exercise, error) {
	const op = "Exercise.Service.CreateExercise"
	log := s.log.With("op", op)
	data, err := s.repo.CreateExercise(ctx, title, desc)
	if err != nil {
		log.Error(op, "Failed to create exercise")
		return nil, err
	}
	return &domain.Exercise{ID: data.ID, Title: data.Title, Description: data.Description}, nil
}
