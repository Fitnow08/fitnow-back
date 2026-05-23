package domain

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID        uuid.UUID  `json:"id"`
	TrainID   uuid.UUID  `json:"train_id"`
	UserId    uuid.UUID  `json:"user_id"`
	ParentID  *uuid.UUID `json:"parent_id"`
	Comment   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	IsDeleted bool       `json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at"`
	Version   int64      `json:"version"`
}
