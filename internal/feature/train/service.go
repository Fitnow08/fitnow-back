package train

import (
	"context"
	"io"
	"log/slog"
	"time"

	trainv1 "github.com/Fitnow08/fitnow-proto/pkg/gen/go/v1/train"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/Sanchir01/fitnow/pkg/db/connect"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TrainRepository interface {
	GetAllPublicTrains(ctx context.Context, param AllTrainsParams) ([]*TrainDB, error)
	GetTrainByID(ctx context.Context, id uuid.UUID) (*TrainDB, error)
	CreateTrain(ctx context.Context, req CreateTrainRequest, userID uuid.UUID) (*TrainDB, error)
	UpdateTrain(ctx context.Context, id uuid.UUID, req UpdateTrainRequest, userID uuid.UUID) (*TrainDB, error)
	DeleteTrain(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetUserTrains(ctx context.Context, userID uuid.UUID) ([]*TrainDB, error)
	AddUserTrain(ctx context.Context, userID, trainID uuid.UUID) error
	RemoveUserTrain(ctx context.Context, userID, trainID uuid.UUID) error
	UpdateTrainImageUrl(ctx context.Context, trainID uuid.UUID, url string) error
	AddTrainExercises(ctx context.Context, trainID uuid.UUID, exercises []TrainExerciseInput) error
	GetTrainExercises(ctx context.Context, trainID uuid.UUID) ([]*TrainExerciseDB, error)
}
type ITrainClient interface {
	GetTrainAndExercises(ctx context.Context, req *trainv1.GetTrainAndExercisesRequest) (*trainv1.GetTrainAndExercisesResponse, error)
	AddTrainExercises(ctx context.Context, req *trainv1.AddTrainExercisesRequest) (*trainv1.AddTrainExercisesResponse, error)
	GetAllAdminTrains(ctx context.Context, req *trainv1.GetAllAdminTrainsRequest) (*trainv1.GetAllAdminTrainsResponse, error)
	GetAllTrains(ctx context.Context, req *trainv1.GetAllTrainsRequest) (*trainv1.GetAllTrainsResponse, error)
	GetTrainById(ctx context.Context, req *trainv1.GetTrainByIdRequest) (*trainv1.GetTrainByIdResponse, error)
	CreateTrain(ctx context.Context, req *trainv1.CreateTrainRequest) (*trainv1.CreateTrainResponse, error)
	UpdateTrain(ctx context.Context, req *trainv1.UpdateTrainRequest) (*trainv1.UpdateTrainResponse, error)
	DeleteTrain(ctx context.Context, req *trainv1.DeleteTrainRequest) (*trainv1.DeleteTrainResponse, error)
	AddTrainImage(ctx context.Context, req *trainv1.AddTrainImageRequest) (*trainv1.AddTrainImageResponse, error)
}
type Service struct {
	log             *slog.Logger
	trainRepository TrainRepository
	client          ITrainClient
	s3              connect.MiniS3Interface
}

func NewService(log *slog.Logger, s3 connect.MiniS3Interface, trainRepository TrainRepository, client ITrainClient) *Service {
	return &Service{log: log, trainRepository: trainRepository, client: client, s3: s3}
}

func (s *Service) GetAllPublicTrains(ctx context.Context, param AllTrainsParams) ([]*domain.Train, error) {
	const op = "Train.Service.GetAllPublicTrains"
	log := s.log.With(slog.String("op", op))

	resp, err := s.client.GetAllTrains(ctx, &trainv1.GetAllTrainsRequest{
		Limit:      int32(param.Limit),
		CategoryId: uuidToString(param.CategoryId),
		Cursor:     param.Cursor,
		Text:       param.Text,
	})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	trains := make([]*domain.Train, 0, len(resp.Trains))
	for _, t := range resp.Trains {
		trains = append(trains, protoToTrain(t))
	}
	return trains, nil
}

func (s *Service) GetTrainByID(ctx context.Context, id uuid.UUID) (*domain.Train, error) {
	const op = "Train.Service.GetTrainByID"
	log := s.log.With(slog.String("op", op))

	resp, err := s.client.GetTrainById(ctx, &trainv1.GetTrainByIdRequest{Id: id.String()})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return protoToTrain(resp.Train), nil
}

func (s *Service) CreateTrain(ctx context.Context, req CreateTrainRequest, userID uuid.UUID) (*domain.Train, error) {
	const op = "Train.Service.CreateTrain"
	log := s.log.With(slog.String("op", op))

	resp, err := s.client.CreateTrain(ctx, &trainv1.CreateTrainRequest{
		Title:      req.Title,
		Type:       typeToProto(req.Type),
		Duration:   req.Duration,
		IsPublic:   req.IsPublic,
		Difficulty: difficultyToProto(req.Difficulty),
		Calories:   req.Calories,
		CategoryId: uuidToString(req.CategoryId),
		UserId:     userID.String(),
	})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return protoToTrain(resp.Train), nil
}

func (s *Service) UpdateTrain(ctx context.Context, id uuid.UUID, req UpdateTrainRequest, userID uuid.UUID) (*domain.Train, error) {
	const op = "Train.Service.UpdateTrain"
	log := s.log.With(slog.String("op", op))

	pbReq := &trainv1.UpdateTrainRequest{
		TrainId: id.String(),
		Title:   req.Title,
		UserId:  userID.String(),
	}
	if req.IsPublic != nil {
		pbReq.IsPublic = *req.IsPublic
	}
	if req.Type != "" {
		t := typeToProto(req.Type)
		pbReq.Type = &t
	}
	if req.Difficulty != "" {
		d := difficultyToProto(req.Difficulty)
		pbReq.Difficulty = &d
	}
	if req.Duration != 0 {
		pbReq.Duration = &req.Duration
	}
	if req.Calories != 0 {
		pbReq.Calories = &req.Calories
	}

	if _, err := s.client.UpdateTrain(ctx, pbReq); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return s.GetTrainByID(ctx, id)
}

func (s *Service) DeleteTrain(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	const op = "Train.Service.DeleteTrain"
	log := s.log.With(slog.String("op", op))

	if _, err := s.client.DeleteTrain(ctx, &trainv1.DeleteTrainRequest{
		Id:     id.String(),
		UserId: userID.String(),
	}); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (s *Service) GetUserTrains(ctx context.Context, userID uuid.UUID) ([]*domain.Train, error) {
	rows, err := s.trainRepository.GetUserTrains(ctx, userID)
	if err != nil {
		return nil, err
	}
	trains := make([]*domain.Train, 0, len(rows))
	for _, t := range rows {
		trains = append(trains, s.dbToDomain(t))
	}
	return trains, nil
}

func (s *Service) AddUserTrain(ctx context.Context, userID, trainID uuid.UUID) error {
	return s.trainRepository.AddUserTrain(ctx, userID, trainID)
}

func (s *Service) RemoveUserTrain(ctx context.Context, userID, trainID uuid.UUID) error {
	return s.trainRepository.RemoveUserTrain(ctx, userID, trainID)
}

func (s *Service) UploadTrainImage(ctx context.Context, trainID uuid.UUID, ext, contentType string, size int64, r io.Reader) error {
	const op = "Train.Service.UploadTrainImage"
	log := s.log.With(slog.String("op", op))

	image, err := io.ReadAll(r)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if _, err := s.client.AddTrainImage(ctx, &trainv1.AddTrainImageRequest{
		TrainId:     trainID.String(),
		Extension:   ext,
		ContentType: contentType,
		Image:       image,
	}); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (s *Service) UploadAllTrainExercises(ctx context.Context, trainID uuid.UUID, exercises []TrainExerciseInput) error {
	const op = "Train.Service.UploadAllTrainExercises"
	log := s.log.With("op", op)

	pbExercises := make([]*trainv1.TrainExerciseInput, 0, len(exercises))
	for _, e := range exercises {
		pbExercises = append(pbExercises, &trainv1.TrainExerciseInput{
			ExerciseId: e.ExerciseID.String(),
			Steps:      int32(e.Steps),
			Sets:       int32(e.Sets),
			Position:   int32(e.Position),
		})
	}

	if _, err := s.client.AddTrainExercises(ctx, &trainv1.AddTrainExercisesRequest{
		TrainId:   trainID.String(),
		Exercises: pbExercises,
	}); err != nil {
		log.Error("fail to add train exercises", "err", err)
		return err
	}
	return nil
}

func (s *Service) GetTrainExercises(ctx context.Context, trainID uuid.UUID) (*TrainAndExercises, error) {
	const op = "Train.Service.GetTrainExercises"
	log := s.log.With("op", op)

	resp, err := s.client.GetTrainAndExercises(ctx, &trainv1.GetTrainAndExercisesRequest{
		StrainId: trainID.String(),
	})
	if err != nil {
		log.Error("failed to get train data", "err", err.Error())
		return nil, err
	}

	exercises := make([]domain.TrainExercise, 0, len(resp.Exercises))
	for _, e := range resp.Exercises {
		exercises = append(exercises, domain.TrainExercise{
			ID:          parseUUID(e.Id),
			Title:       e.Title,
			Description: e.Description,
			VideoURL:    e.VideoUrl,
			Difficulty:  difficultyToString(e.Difficulty),
			Steps:       int(e.Steps),
			Sets:        int(e.Sets),
			Position:    int(e.Position),
		})
	}

	return &TrainAndExercises{
		Train:     *protoToTrain(resp.Train),
		Exercises: exercises,
	}, nil
}

func (s *Service) dbToDomain(t *TrainDB) *domain.Train {
	return &domain.Train{
		ID:         t.ID,
		Title:      t.Title,
		Type:       t.Type,
		Duration:   t.Duration,
		IsPublic:   t.IsPublic,
		Difficulty: t.Difficulty,
		CategoryId: t.CategoryId,
		Calories:   t.Calories,
		ImageURL:   s.s3.PublicURL(t.ImagePath),
		CreatedBy:  t.CreatedBy,
		CreatedAt:  t.CreatedAt,
	}
}

func protoToTrain(t *trainv1.Train) *domain.Train {
	if t == nil {
		return &domain.Train{}
	}
	return &domain.Train{
		ID:         parseUUID(t.Id),
		Title:      t.Title,
		Type:       t.Type,
		Duration:   t.Duration,
		IsPublic:   t.IsPublic,
		Difficulty: t.Difficulty,
		Calories:   t.Calories,
		CategoryId: parseUUID(t.CategoryId),
		ImageURL:   t.ImagePath,
		CreatedBy:  parseUUID(t.CreatedBy),
		CreatedAt:  tsToTime(t.CreatedAt),
	}
}

func parseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func uuidToString(id uuid.UUID) string {
	if id == uuid.Nil {
		return ""
	}
	return id.String()
}

func tsToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func typeToProto(t string) trainv1.TrainType {
	switch t {
	case "strength":
		return trainv1.TrainType_TRAIN_TYPE_STRENGTH
	case "cardio":
		return trainv1.TrainType_TRAIN_TYPE_CARDIO
	case "stretching":
		return trainv1.TrainType_TRAIN_TYPE_STRETCHING
	default:
		return trainv1.TrainType_TRAIN_TYPE_UNSPECIFIED
	}
}

func difficultyToProto(d string) trainv1.Difficulty {
	switch d {
	case "easy":
		return trainv1.Difficulty_DIFFICULTY_EASY
	case "medium":
		return trainv1.Difficulty_DIFFICULTY_MEDIUM
	case "hard":
		return trainv1.Difficulty_DIFFICULTY_HARD
	default:
		return trainv1.Difficulty_DIFFICULTY_UNSPECIFIED
	}
}

func difficultyToString(d trainv1.Difficulty) string {
	switch d {
	case trainv1.Difficulty_DIFFICULTY_EASY:
		return "easy"
	case trainv1.Difficulty_DIFFICULTY_MEDIUM:
		return "medium"
	case trainv1.Difficulty_DIFFICULTY_HARD:
		return "hard"
	default:
		return ""
	}
}
