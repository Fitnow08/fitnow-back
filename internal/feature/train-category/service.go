package traincategory

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
)

type TrainCategoryRepoInterface interface {
	DeleteTrainCategory(ctx context.Context, id uuid.UUID) error
	UpdateTrainCategory(ctx context.Context, id uuid.UUID, title string) (*TrainCategoryDB, error)
	CreateTrainCategory(ctx context.Context, title string) (*TrainCategoryDB, error)
	GetAllTrainCategory(ctx context.Context) ([]TrainCategoryDB, error)
}
type Service struct {
	log    *slog.Logger
	tcrepo TrainCategoryRepoInterface
}

func NewService(log *slog.Logger, tcrepo TrainCategoryRepoInterface) *Service {
	return &Service{
		log:    log,
		tcrepo: tcrepo,
	}
}

func (s *Service) GetAllTrainCategory(ctx context.Context) ([]TrainCategoryDB, error) {
	const op = "GetAllTrain.Service.GetAllTrainCategory"
	log := s.log.With(slog.String("op", op))
	category, err := s.tcrepo.GetAllTrainCategory(ctx)
	if err != nil {
		log.Error("fail to get all train category")
		return nil, err
	}
	log.Info("success to get all train category")
	return category, nil
}

func (s *Service) CreateTrainCategory(ctx context.Context, title string) (*TrainCategoryDB, error) {
	const op = "TrainCategory.Service.CreateTrainCategory"
	log := s.log.With(slog.String("op", op))

	category, err := s.tcrepo.CreateTrainCategory(ctx, title)
	if err != nil {
		log.Error("fail to create train category")
		return nil, err
	}
	log.Info("success to create train category")
	return category, nil
}
func (s *Service) UpdateTrainCategory(ctx context.Context, id uuid.UUID, title string) (*TrainCategoryDB, error) {
	const op = "TrainCategory.Service.UpdateTrainCategory"
	log := s.log.With(slog.String("op", op))

	category, err := s.tcrepo.UpdateTrainCategory(ctx, id, title)
	if err != nil {
		log.Error("fail to update train category")
		return nil, err
	}
	log.Info("success to update train category")
	return category, nil
}
func (s *Service) DeleteTrainCategory(ctx context.Context, id uuid.UUID) error {
	const op = "TrainCategory.Service.DeleteTrainCategory"
	log := s.log.With(slog.String("op", op))

	if err := s.tcrepo.DeleteTrainCategory(ctx, id); err != nil {
		log.Error("fail to delete train category")
		return err
	}
	log.Info("success to delete train category")
	return nil
}
