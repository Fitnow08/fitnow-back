package programcategory

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
)

type ProgramCategoryRepoInterface interface {
	GetAllProgramCategory(ctx context.Context) ([]ProgramCategoryDB, error)
	CreateProgramCategory(ctx context.Context, title string) (*ProgramCategoryDB, error)
	UpdateProgramCategory(ctx context.Context, id uuid.UUID, title string) (*ProgramCategoryDB, error)
	DeleteProgramCategory(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	log    *slog.Logger
	pcrepo ProgramCategoryRepoInterface
}

func NewService(log *slog.Logger, pcrepo ProgramCategoryRepoInterface) *Service {
	return &Service{
		log:    log,
		pcrepo: pcrepo,
	}
}

func (s *Service) GetAllProgramCategory(ctx context.Context) ([]ProgramCategoryDB, error) {
	const op = "ProgramCategory.Service.GetAllProgramCategory"
	log := s.log.With(slog.String("op", op))
	category, err := s.pcrepo.GetAllProgramCategory(ctx)
	if err != nil {
		log.Error("fail to get all program category")
		return nil, err
	}
	log.Info("success to get all program category")
	return category, nil
}

func (s *Service) CreateProgramCategory(ctx context.Context, title string) (*ProgramCategoryDB, error) {
	const op = "ProgramCategory.Service.CreateProgramCategory"
	log := s.log.With(slog.String("op", op))

	category, err := s.pcrepo.CreateProgramCategory(ctx, title)
	if err != nil {
		log.Error("fail to create program category")
		return nil, err
	}
	log.Info("success to create program category")
	return category, nil
}

func (s *Service) UpdateProgramCategory(ctx context.Context, id uuid.UUID, title string) (*ProgramCategoryDB, error) {
	const op = "ProgramCategory.Service.UpdateProgramCategory"
	log := s.log.With(slog.String("op", op))

	category, err := s.pcrepo.UpdateProgramCategory(ctx, id, title)
	if err != nil {
		log.Error("fail to update program category")
		return nil, err
	}
	log.Info("success to update program category")
	return category, nil
}

func (s *Service) DeleteProgramCategory(ctx context.Context, id uuid.UUID) error {
	const op = "ProgramCategory.Service.DeleteProgramCategory"
	log := s.log.With(slog.String("op", op))

	if err := s.pcrepo.DeleteProgramCategory(ctx, id); err != nil {
		log.Error("fail to delete program category")
		return err
	}
	log.Info("success to delete program category")
	return nil
}
