package comment

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/google/uuid"
	"log/slog"
)

type RepositoryInterface interface {
	DeleteComment(ctx context.Context, commentID uuid.UUID) error
	UpdateComment(ctx context.Context, comment string, commentID uuid.UUID) error
	GetTrainComments(ctx context.Context, trainID uuid.UUID) ([]CommentDB, error)
	CreateComment(ctx context.Context, comment string, train_id, user_id uuid.UUID, parentid *uuid.UUID) error
}
type Service struct {
	repo RepositoryInterface
	log  *slog.Logger
}

func NewService(log *slog.Logger, repo RepositoryInterface) *Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) GetTrainComments(ctx context.Context, trainID uuid.UUID) ([]domain.Comment, error) {
	const op = "Comment.Service.GetTrainComments"
	log := s.log.With(slog.String("op", op))
	rows, err := s.repo.GetTrainComments(ctx, trainID)
	if err != nil {
		log.Error("failed get comments train", slog.Any("err", err))
		return nil, err
	}
	comments := make([]domain.Comment, 0, len(rows))
	for _, c := range rows {
		comments = append(comments, dbToDomain(c))
	}
	return comments, nil
}

func (s *Service) CreateComment(ctx context.Context, comment string, train_id, user_id uuid.UUID, parentid *uuid.UUID) error {
	const op = "Comment.Service.CreateComment"
	log := s.log.With(slog.String("op", op))
	err := s.repo.CreateComment(ctx, comment, train_id, user_id, parentid)
	if err != nil {
		log.Error("failed create comment", err.Error())
		return err
	}
	return nil
}

func (s *Service) UpdateComment(ctx context.Context, comment string, commentID uuid.UUID) error {
	const op = "Comment.Service.UpdateComment"
	log := s.log.With(slog.String("op", op))
	err := s.repo.UpdateComment(ctx, comment, commentID)
	if err != nil {
		log.Error("failed update comment", err.Error())
		return err
	}
	return nil
}
func (s *Service) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	const op = "Comment.Service.DeleteComment"
	log := s.log.With(slog.String("op", op))
	err := s.repo.DeleteComment(ctx, commentID)
	if err != nil {
		log.Error("failed delete comment", err.Error())
		return err
	}
	return nil
}
