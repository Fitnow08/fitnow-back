package rating

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
)

type RatingRepository interface {
	GetAllRatings(ctx context.Context) ([]RatingDB, error)
	CreateTrainRating(ctx context.Context, userid, trainid uuid.UUID, rating int) error
	UpdateTrainRating(ctx context.Context, userid, trainid uuid.UUID, rating int) error
}

type Service struct {
	log  *slog.Logger
	repo RatingRepository
}

func NewService(log *slog.Logger, repo RatingRepository) *Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) GetAllTrainRatings(ctx context.Context) ([]RatingDB, error) {
	const op = "Rating.Service.GetAllRatings"
	log := s.log.With("op", op)
	ratings, err := s.repo.GetAllRatings(ctx)
	if err != nil {
		log.Error("failed get all rating", err.Error())
		return nil, err
	}

	return ratings, nil
}

func (s *Service) CreateTrainRating(ctx context.Context, userid, trainid uuid.UUID, rating int) error {
	const op = "Rating.Service.CreateTrainRating"
	log := s.log.With("op", op)
	if err := s.repo.CreateTrainRating(ctx, userid, trainid, rating); err != nil {
		log.Error("failed create train rating", err.Error())
		return err
	}
	return nil
}

func (s *Service) UpdateTrainRating(ctx context.Context, userid, trainid uuid.UUID, rating int) error {
	const op = "Rating.Service.UpdateTrainRating"
	log := s.log.With("op", op)
	if err := s.repo.UpdateTrainRating(ctx, userid, trainid, rating); err != nil {
		log.Error("failed update train rating", err.Error())
		return err
	}
	return nil
}
