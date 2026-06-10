package events

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	constants "github.com/Sanchir01/fitnow/internal/models/contants"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

type Repository struct {
	log *slog.Logger
	db  *pgxpool.Pool
}

func NewRepository(log *slog.Logger, db *pgxpool.Pool) *Repository {
	return &Repository{log: log, db: db}
}

func (r *Repository) GetNewEvents(ctx context.Context, limit uint64) ([]domain.Event, error) {
	query, args, err := sq.Select("id", "event_type", "payload").
		From(constants.EventsTableName).
		Where(sq.Eq{"event_type": "new"}).
		Limit(limit).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		if err := rows.Scan(&event.ID, &event.EventType, &event.Payload); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *Repository) CreateEvent(ctx context.Context, eventype, payload string) error {
	query, args, err := sq.
		Insert(constants.EventsTableName).
		Columns("event_type", "payload", "reserved_to").
		Values(eventype, payload, nil).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) ClaimEvents(ctx context.Context, limit int, lease time.Duration) ([]domain.Event, error) {
	const q = `
  UPDATE events
     SET reserved_to = $1
   WHERE id IN (
     SELECT id FROM events
      WHERE status = 'new'
        AND (reserved_to IS NULL OR reserved_to < now())
      ORDER BY created_at
      LIMIT $2
      FOR UPDATE SKIP LOCKED
   )
  RETURNING id, event_type, payload;`

	rows, err := r.db.Query(ctx, q, time.Now().Add(lease), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(&e.ID, &e.EventType, &e.Payload); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *Repository) UpdateEvent(ctx context.Context, eventsids []uuid.UUID) error {
	if len(eventsids) == 0 {
		return nil
	}
	query, args, err := sq.
		Update(constants.EventsTableName).
		Set("status", "done").
		Where(sq.Eq{"id": eventsids}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
