package comment

import (
	"github.com/google/uuid"
	"time"
)

type CommentDB struct {
	ID        uuid.UUID  `db:"id"`
	TrainID   uuid.UUID  `db:"train_id"`
	UserId    uuid.UUID  `db:"user_id"`
	ParentID  *uuid.UUID `db:"parent_id"`
	Comment   string     `db:"content"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	IsDeleted bool       `db:"is_deleted"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int64      `db:"version"`
}
