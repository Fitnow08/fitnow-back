package comment

import (
	"github.com/Sanchir01/fitnow/internal/models/domain"
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

type CreateCommentRequest struct {
	Comment  string     `json:"comment" validate:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
}

type UpdateCommentRequest struct {
	Comment  string     `json:"comment" validate:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
}

func dbToDomain(c CommentDB) domain.Comment {
	return domain.Comment{
		ID:        c.ID,
		TrainID:   c.TrainID,
		UserId:    c.UserId,
		ParentID:  c.ParentID,
		Comment:   c.Comment,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		IsDeleted: c.IsDeleted,
		DeletedAt: c.DeletedAt,
		Version:   c.Version,
	}
}
