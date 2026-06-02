package exercises

import (
	"context"
	"fmt"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/Sanchir01/fitnow/pkg/db/connect"
	"github.com/google/uuid"
	"io"
	"log/slog"
)

type ExercisesRepository interface {
	GetAllExercises(ctx context.Context) ([]ExerciseDB, error)
	CreateExercise(ctx context.Context, title, desc string) (*ExerciseDB, error)
	UpdateVideoUrl(ctx context.Context, url string, id uuid.UUID) error
}

type Service struct {
	log  *slog.Logger
	repo ExercisesRepository
	s3   connect.MiniS3Interface
}

func NewService(log *slog.Logger, repo ExercisesRepository, s3 connect.MiniS3Interface) *Service {
	return &Service{log: log, repo: repo, s3: s3}
}

func (s *Service) GetAllExercises(ctx context.Context) ([]domain.Exercise, error) {
	const op = "Exercise.Service.GetAllExercises"
	log := s.log.With("op", op)
	data, err := s.repo.GetAllExercises(ctx)
	if err != nil {
		log.Error("failed to get all exercises", slog.Any("err", err.Error()))
		return nil, err
	}
	var exercises = make([]domain.Exercise, 0, len(data))
	for _, info := range data {
		exercises = append(exercises, domain.Exercise{
			ID:          info.ID,
			Title:       info.Title,
			Description: info.Description,
			VideoURL:    info.VideoURL,
		})
	}
	return exercises, nil
}

func (s *Service) CreateExercise(ctx context.Context, title, desc string) (*domain.Exercise, error) {
	const op = "Exercise.Service.CreateExercise"
	log := s.log.With("op", op)
	data, err := s.repo.CreateExercise(ctx, title, desc)
	if err != nil {
		log.Error("failed to create exercise", slog.Any("err", err))
		return nil, err
	}
	return &domain.Exercise{ID: data.ID, Title: data.Title, Description: data.Description}, nil
}

func (s *Service) SaveExerciseVideo(ctx context.Context, exercise uuid.UUID, ext, contentType string, size int64, r io.Reader) error {
	const op = "Exercise.Service.SaveExerciseVideo"
	log := s.log.With("op", op)
	key := fmt.Sprintf("exercises/%s%s", exercise, ext)
	if err := s.s3.Upload(ctx, key, r, size, contentType); err != nil {
		log.Error(err.Error())
		return err
	}
	if err := s.repo.UpdateVideoUrl(ctx, key, exercise); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil

}
