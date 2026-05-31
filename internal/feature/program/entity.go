package program

import (
	"github.com/google/uuid"
	"time"
)

type CreateProgramRequest struct {
	Title       string     `json:"title" validate:"required,max=255"`
	Description string     `json:"description" validate:"required,max=255"`
	Weeks       int        `json:"weeks" validate:"required"`
	Difficulty  Level      `json:"difficulty" validate:"required"`
	CategoryID  *uuid.UUID `json:"category_id" validate:"omitempty"`
}
type ProgramDB struct {
	ID         uuid.UUID  `db:"id"`
	Title      string     `db:"title"`
	Desc       string     `db:"description"`
	Weeks      int        `db:"weeks"`
	Difficult  string     `db:"difficulty"`
	IsPublic   bool       `db:"is_public"`
	CategoryID *uuid.UUID `db:"category_id"`
	ImagePath  string     `db:"image_path"`
	CreatedBy  uuid.UUID  `db:"created_by"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	Version    int64      `db:"version"`
}

type ProgramAndCountTrainDB struct {
	ProgramDB
	TrainsCount int64 `db:"trains_count"`
}
