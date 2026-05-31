package program

import (
	"context"
	"fmt"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/Sanchir01/fitnow/pkg/db/connect"
	"github.com/google/uuid"
	"io"
	"log/slog"
)

type ProgramRepository interface {
	CreateProgram(ctx context.Context, title string, description string, weeks int, level Level, categoryID *uuid.UUID, user_id uuid.UUID) (*ProgramDB, error)
	GetAllPrograms(ctx context.Context) ([]ProgramAndCountTrainDB, error)
	UpdateProgramImagePath(ctx context.Context, programID uuid.UUID, imagePath string) error
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

func (s *Service) GetAllProgramAndTrainsCount(ctx context.Context) ([]domain.ProgramAndTrainsCount, error) {
	const op = "Program.Service.GetAllProgramAndTrainsCount"
	log := s.log.With(slog.String("op", op))
	programs, err := s.repo.GetAllPrograms(ctx)
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
