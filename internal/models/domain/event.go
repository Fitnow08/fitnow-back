package domain

import "github.com/google/uuid"

type Event struct {
	ID        uuid.UUID `json:"id"`
	EventType string    `json:"event_type"`
	Payload   string    `json:"payload"`
}
