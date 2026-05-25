package app

import (
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"github.com/Sanchir01/fitnow/internal/feature/comment"
	"github.com/Sanchir01/fitnow/internal/feature/exercises"
	"github.com/Sanchir01/fitnow/internal/feature/rating"
	"github.com/Sanchir01/fitnow/internal/feature/train"
	traincategory "github.com/Sanchir01/fitnow/internal/feature/train-category"
	"log/slog"
)

type Services struct {
	AuthService          *auth.Service
	TrainService         *train.Service
	ExercisesService     *exercises.Service
	TrainCategoryService *traincategory.Service
	RatingService        *rating.Service
	CommentService       *comment.Service
}

func NewServices(repo *Repositories, s3minio *S3, l *slog.Logger) *Services {
	return &Services{
		AuthService:          auth.NewService(l, repo.AuthRepository),
		TrainService:         train.NewService(l, s3minio.TrainBucket, repo.TrainRepository),
		ExercisesService:     exercises.NewService(l, repo.ExercisesRepository),
		TrainCategoryService: traincategory.NewService(l, repo.TrainCategoryRepository),
		RatingService:        rating.NewService(l, repo.RatingRepository),
		CommentService:       comment.NewService(l, repo.CommentRepository),
	}
}
