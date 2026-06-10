package events

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/Sanchir01/fitnow/pkg/kafka"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, eventype, payload string) error
	GetNewEvents(ctx context.Context, limit uint64) ([]domain.Event, error)
	UpdateEvent(ctx context.Context, eventsids []uuid.UUID) error
	ClaimEvents(ctx context.Context, limit int, lease time.Duration) ([]domain.Event, error)
}

type Service struct {
	log          *slog.Logger
	repo         EventRepository
	kafkaservice kafka.EventProducer
}

func NewService(log *slog.Logger, repo EventRepository, kafkaservice kafka.EventProducer) *Service {
	return &Service{log: log, repo: repo, kafkaservice: kafkaservice}
}

func (s *Service) StartProcessEvents(ctx context.Context, handreperiod time.Duration) {
	const op = "Events.Service.StartProcessEvents"
	log := s.log.With("op", op)
	ticker := time.NewTicker(handreperiod)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Info("shutting down events service")
				return
			case <-ticker.C:

			}
			events, err := s.repo.ClaimEvents(ctx, 100, 5*time.Minute)
			if err != nil {
				log.Error("failed to get new events", err.Error())
				continue
			}
			if len(events) == 0 {
				log.Debug("no new events found")
				continue
			}

			sentIDs := make([]uuid.UUID, 0, len(events))
			for _, event := range events {
				if err := s.kafkaservice.Publish(
					ctx,
					event.EventType,
					[]byte(event.ID.String()),
					[]byte(event.Payload),
				); err != nil {
					log.Error("failed to publish event", "id", event.ID, "err", err.Error())
					break
				}
				sentIDs = append(sentIDs, event.ID)
			}

			if len(sentIDs) == 0 {
				continue
			}

			if len(sentIDs) == 0 {
				continue
			}

			if err := s.repo.UpdateEvent(ctx, sentIDs); err != nil {
				log.Error("failed to update events", err.Error())
				continue
			}
		}
	}()
}
