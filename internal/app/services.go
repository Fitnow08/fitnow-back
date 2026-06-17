package app

import (
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"github.com/Sanchir01/fitnow/internal/feature/comment"
	"github.com/Sanchir01/fitnow/internal/feature/events"
	"github.com/Sanchir01/fitnow/internal/feature/exercises"
	"github.com/Sanchir01/fitnow/internal/feature/program"
	programcategory "github.com/Sanchir01/fitnow/internal/feature/program-category"
	"github.com/Sanchir01/fitnow/internal/feature/rating"
	"github.com/Sanchir01/fitnow/internal/feature/train"
	traincategory "github.com/Sanchir01/fitnow/internal/feature/train-category"
	"github.com/Sanchir01/fitnow/internal/feature/user"
	"github.com/Sanchir01/fitnow/pkg/kafka"
	"log/slog"
)

type Services struct {
	AuthService            *auth.Service
	TrainService           *train.Service
	ExercisesService       *exercises.Service
	TrainCategoryService   *traincategory.Service
	RatingService          *rating.Service
	CommentService         *comment.Service
	ProgramService         *program.Service
	ProgramCategoryService *programcategory.Service
	UserService            *user.Service
	EventService           *events.Service
}

func NewServices(repo *Repositories, s3minio *S3, clients *Clients, producer *kafka.Producer, l *slog.Logger) *Services {

	return &Services{
		AuthService:            auth.NewService(l, clients.AuthClient),
		TrainService:           train.NewService(l, s3minio.TrainBucket, repo.TrainRepository, clients.TrainClient),
		ExercisesService:       exercises.NewService(l, repo.ExercisesRepository, s3minio.ExercisesBucket),
		TrainCategoryService:   traincategory.NewService(l, repo.TrainCategoryRepository),
		RatingService:          rating.NewService(l, repo.RatingRepository),
		CommentService:         comment.NewService(l, repo.CommentRepository),
		ProgramService:         program.NewService(l, clients.ProgramClient, repo.EventRepository),
		ProgramCategoryService: programcategory.NewService(l, repo.ProgramCategoryRepository),
		UserService:            user.NewService(l, clients.AuthClient),
		EventService:           events.NewService(l, repo.EventRepository, producer),
	}
}
