package program

import (
	"context"
	"encoding/json"
	"github.com/Sanchir01/fitnow/internal/feature/events"
	constants "github.com/Sanchir01/fitnow/internal/models/contants"
	"io"
	"log/slog"
	"time"

	programv1 "github.com/Fitnow08/fitnow-proto/pkg/gen/go/v1/program"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProgramClient interface {
	GetAllPrograms(ctx context.Context, categoryID, search string) (*programv1.GetAllProgramsResponse, error)
	CreateProgram(ctx context.Context, req *programv1.CreateProgramRequest) (*programv1.CreateProgramResponse, error)
	AddProgramImage(ctx context.Context, req *programv1.AddProgramImageRequest) (*programv1.AddProgramImageResponse, error)
	GetProgramsAndTrains(ctx context.Context, programID string) (*programv1.GetProgramsAndTrainsResponse, error)
	UploadAllProgramTrains(ctx context.Context, programID string, trains []*programv1.ProgramTrainInput) (*programv1.UploadAllProgramTrainsResponse, error)
}

type Service struct {
	log          *slog.Logger
	client       ProgramClient
	eventservice events.EventRepository
}

func NewService(log *slog.Logger, client ProgramClient, eventservice events.EventRepository) *Service {
	return &Service{
		log:          log,
		client:       client,
		eventservice: eventservice,
	}
}

func (s *Service) GetAllProgramAndTrainsCount(ctx context.Context, categoryId uuid.UUID, search string) ([]domain.ProgramAndTrainsCount, error) {
	const op = "Program.Service.GetAllProgramAndTrainsCount"
	log := s.log.With(slog.String("op", op))

	var categoryID string
	if categoryId != uuid.Nil {
		categoryID = categoryId.String()
	}

	resp, err := s.client.GetAllPrograms(ctx, categoryID, search)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	programs := make([]domain.ProgramAndTrainsCount, 0, len(resp.Programs))
	for _, p := range resp.Programs {
		programs = append(programs, domain.ProgramAndTrainsCount{
			ID:          parseUUID(p.Id),
			Title:       p.Title,
			Desc:        p.Desc,
			Weeks:       int(p.Weeks),
			Difficult:   p.Difficulty,
			IsPublic:    p.IsPublic,
			CategoryID:  parseUUIDPtr(p.CategoryId),
			ImagePath:   p.ImagePath,
			ImageURL:    p.ImageUrl,
			CreatedBy:   parseUUID(p.CreatedBy),
			CreatedAt:   tsToTime(p.CreatedAt),
			UpdatedAt:   tsToTime(p.UpdatedAt),
			Version:     p.Version,
			TrainsCount: p.TrainsCount,
		})
	}
	return programs, nil
}

func (s *Service) CreateProgram(ctx context.Context, title string, description string, weeks int, level domain.Level, categoryID *uuid.UUID, user_id uuid.UUID) (*domain.Program, error) {
	const op = "Program.Service.CreateProgram"
	log := s.log.With(slog.String("op", op))

	var categoryid *string
	if categoryID != nil && *categoryID != uuid.Nil {
		id := categoryID.String()
		categoryid = &id
	}

	if _, err := s.client.CreateProgram(ctx, &programv1.CreateProgramRequest{
		Title:       title,
		Description: description,
		Weeks:       int32(weeks),
		Difficulty:  LevelToProto(level),
		CategoryId:  categoryid,
		UserId:      user_id.String(),
	}); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &domain.Program{
		Title:      title,
		Desc:       description,
		Weeks:      weeks,
		Difficult:  string(level),
		CategoryID: categoryID,
		CreatedBy:  user_id,
	}, nil
}

func (s *Service) PublishProgramCreated(ctx context.Context, userID uuid.UUID, req CreateProgramRequest) error {
	const op = "Program.Service.PublishProgramCreated"
	log := s.log.With(slog.String("op", op))

	event := domain.ProgramCreatedEvent{
		Title:       req.Title,
		Description: req.Description,
		Weeks:       req.Weeks,
		Difficulty:  req.Difficulty,
		CategoryID:  req.CategoryID,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
	}

	value, err := json.Marshal(event)
	if err != nil {
		log.Error("failed to marshal program.created event", slog.String("err", err.Error()))
		return err
	}
	if err := s.eventservice.CreateEvent(ctx, constants.TopicProgramCreated, string(value)); err != nil {
		log.Error("failed to publish program.created event", slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (s *Service) UploadProgramImage(ctx context.Context, programID uuid.UUID, ext, contentType string, size int64, r io.Reader) error {
	const op = "Program.Service.UploadProgramImage"
	log := s.log.With(slog.String("op", op))

	image, err := io.ReadAll(r)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if _, err := s.client.AddProgramImage(ctx, &programv1.AddProgramImageRequest{
		ProgramId:   programID.String(),
		Extension:   ext,
		ContentType: contentType,
		Image:       image,
	}); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (s *Service) UploadAllProgramTrains(ctx context.Context, programId uuid.UUID, trains []ProgramTrainInput) error {
	const op = "Program.Service.UploadAllProgramTrains"
	log := s.log.With(slog.String("op", op))

	pbTrains := make([]*programv1.ProgramTrainInput, 0, len(trains))
	for _, t := range trains {
		pbTrains = append(pbTrains, &programv1.ProgramTrainInput{
			TrainId:    t.TrainId.String(),
			WeekNumber: int32(t.WeekNumber),
			DayOfWeek:  int32(t.DayOfWeek),
			Position:   int32(t.Position),
		})
	}

	if _, err := s.client.UploadAllProgramTrains(ctx, programId.String(), pbTrains); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (s *Service) GetProgramTrains(ctx context.Context, programID uuid.UUID) (*domain.ProgramAndTrain, error) {
	const op = "Program.Service.GetProgramTrains"
	log := s.log.With(slog.String("op", op))

	resp, err := s.client.GetProgramsAndTrains(ctx, programID.String())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	p := resp.Program
	trains := make([]domain.ProgramTrains, 0, len(resp.Trains))
	for _, t := range resp.Trains {
		trains = append(trains, domain.ProgramTrains{
			Train:      protoToTrain(t.Train),
			WeekNumber: int(t.WeekNumber),
			DayOfWeek:  int(t.DayOfWeek),
			Position:   int(t.Position),
		})
	}

	return &domain.ProgramAndTrain{
		Program: domain.Program{
			ID:         parseUUID(p.Id),
			Title:      p.Title,
			Desc:       p.Desc,
			Weeks:      int(p.Weeks),
			Difficult:  p.Difficulty,
			IsPublic:   p.IsPublic,
			CategoryID: parseUUIDPtr(p.CategoryId),
			ImageURL:   p.ImageUrl,
			CreatedBy:  parseUUID(p.UserId),
			CreatedAt:  tsToTime(p.CreatedAt),
			UpdatedAt:  tsToTime(p.UpdatedAt),
			Version:    p.Version,
		},
		Trains: trains,
	}, nil
}

func protoToTrain(t *programv1.Train) domain.Train {
	if t == nil {
		return domain.Train{}
	}
	return domain.Train{
		ID:         parseUUID(t.Id),
		Title:      t.Title,
		Type:       t.Type,
		Duration:   t.Duration,
		IsPublic:   t.IsPublic,
		Difficulty: t.Difficulty,
		Calories:   t.Calories,
		CategoryId: parseUUID(t.CategoryId),
		ImageURL:   t.ImageUrl,
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

func parseUUIDPtr(s *string) *uuid.UUID {
	if s == nil || *s == "" {
		return nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil
	}
	return &id
}

func tsToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
