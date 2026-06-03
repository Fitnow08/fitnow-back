package program

import (
	"context"
	"fmt"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/Sanchir01/fitnow/pkg/db/connect"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"io"
	"log/slog"
)

type ProgramRepository interface {
	CreateProgram(ctx context.Context, title string, description string, weeks int, level Level, categoryID *uuid.UUID, user_id uuid.UUID) (*ProgramDB, error)
	GetAllPrograms(ctx context.Context, categoryId uuid.UUID, search string) ([]ProgramAndCountTrainDB, error)
	UpdateProgramImagePath(ctx context.Context, programID uuid.UUID, imagePath string) error
	AddProgramTrains(ctx context.Context, programId uuid.UUID, trains []ProgramTrainInput) error
	GetProgramTrains(ctx context.Context, programId uuid.UUID) ([]ProgramTrainDB, error)
	GetProgramById(ctx context.Context, programID uuid.UUID) (*ProgramDB, error)
}
type Service struct {
	log  *slog.Logger
	repo ProgramRepository
	s3   connect.MiniS3Interface
}

func NewService(log *slog.Logger, repo ProgramRepository, s3 connect.MiniS3Interface) *Service {
	return &Service{
		log:  log,
		repo: repo,
		s3:   s3,
	}
}

func (s *Service) GetAllProgramAndTrainsCount(ctx context.Context, categoryId uuid.UUID, search string) ([]domain.ProgramAndTrainsCount, error) {
	const op = "Program.Service.GetAllProgramAndTrainsCount"
	log := s.log.With(slog.String("op", op))
	programs, err := s.repo.GetAllPrograms(ctx, categoryId, search)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	finalprograms := make([]domain.ProgramAndTrainsCount, 0, len(programs))
	for _, program := range programs {

		finalprograms = append(finalprograms, domain.ProgramAndTrainsCount{
			ID:          program.ID,
			Title:       program.Title,
			Desc:        program.Desc,
			Weeks:       program.Weeks,
			CreatedBy:   program.CreatedBy,
			CreatedAt:   program.CreatedAt,
			UpdatedAt:   program.UpdatedAt,
			CategoryID:  program.CategoryID,
			ImageURL:    s.s3.PublicURL(program.ImagePath),
			Version:     program.Version,
			TrainsCount: program.TrainsCount,
			IsPublic:    program.IsPublic,
			Difficult:   program.Difficult,
		})
	}
	return finalprograms, nil
}
func (s *Service) CreateProgram(ctx context.Context, title string, description string, weeks int, level Level, categoryID *uuid.UUID, user_id uuid.UUID) (*domain.Program, error) {
	const op = "Program.Service.CreateProgram"
	log := s.log.With(slog.String("op", op))

	program, err := s.repo.CreateProgram(ctx, title, description, weeks, level, categoryID, user_id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &domain.Program{
		ID:         program.ID,
		Title:      program.Title,
		Desc:       program.Desc,
		Weeks:      weeks,
		Difficult:  program.Difficult,
		IsPublic:   program.IsPublic,
		CategoryID: program.CategoryID,
		ImageURL:   s.s3.PublicURL(program.ImagePath),
		CreatedBy:  program.CreatedBy,
		CreatedAt:  program.CreatedAt,
		UpdatedAt:  program.UpdatedAt,
		Version:    program.Version,
	}, nil
}
func (s *Service) UploadProgramImage(ctx context.Context, programID uuid.UUID, ext, contentType string, size int64, r io.Reader) error {
	const op = "Program.Service.UploadProgramImage"
	log := s.log.With(slog.String("op", op))
	key := fmt.Sprintf("programs/%s%s", programID, ext)
	if err := s.s3.Upload(ctx, key, r, size, contentType); err != nil {
		log.Error(err.Error())
		return err
	}
	if err := s.repo.UpdateProgramImagePath(ctx, programID, key); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (s *Service) UploadAllProgramTrains(ctx context.Context, programId uuid.UUID, trains []ProgramTrainInput) error {
	const op = "Program.Service.UploadAllProgramTrains"
	log := s.log.With(slog.String("op", op))
	if err := s.repo.AddProgramTrains(ctx, programId, trains); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (s *Service) GetProgramTrains(ctx context.Context, programID uuid.UUID) (*domain.ProgramAndTrain, error) {
	const op = "Program.Service.GetProgramTrains"
	log := s.log.With(slog.String("op", op))
	g, ctx := errgroup.WithContext(ctx)

	var (
		trains  []ProgramTrainDB
		program *ProgramDB
	)
	g.Go(func() error {
		var err error
		program, err = s.repo.GetProgramById(ctx, programID)
		return err
	})

	g.Go(func() error {
		var err error
		trains, err = s.repo.GetProgramTrains(ctx, programID)
		return err
	})
	if err := g.Wait(); err != nil {
		log.Error("failed to get programs data", "err", err.Error())
		return nil, err
	}
	newtrains := make([]domain.ProgramTrains, 0, len(trains))
	for _, train := range trains {
		t := train.Train
		t.ImageURL = s.s3.PublicURL(train.ImageURL)
		newtrains = append(newtrains, domain.ProgramTrains{
			Train:      t,
			WeekNumber: train.WeekNumber,
			DayOfWeek:  train.DayOfWeek,
			Position:   train.Position,
		})
	}

	var categoryID uuid.UUID
	if program.CategoryID != nil {
		categoryID = *program.CategoryID
	}

	return &domain.ProgramAndTrain{
		Program: domain.ProgramTrains{Train: domain.Train{
			ID:         program.ID,
			Title:      program.Title,
			Difficulty: program.Difficult,
			IsPublic:   program.IsPublic,
			CategoryId: categoryID,
			ImageURL:   s.s3.PublicURL(program.ImagePath),
			CreatedBy:  program.CreatedBy,
			CreatedAt:  program.CreatedAt,
		}},
		Trains: newtrains,
	}, nil
}
