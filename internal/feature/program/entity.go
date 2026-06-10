package program

import (
	programv1 "github.com/Fitnow08/fitnow-proto/pkg/gen/go/v1/program"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/google/uuid"
	"time"
)

func LevelToProto(level domain.Level) programv1.DifficultyLevel {
	switch level {
	case domain.Easy:
		return programv1.DifficultyLevel_DIFFICULTY_LEVEL_EASY
	case domain.Medium:
		return programv1.DifficultyLevel_DIFFICULTY_LEVEL_MEDIUM
	case domain.Hard:
		return programv1.DifficultyLevel_DIFFICULTY_LEVEL_HARD
	default:
		return programv1.DifficultyLevel_DIFFICULTY_LEVEL_UNSPECIFIED
	}
}

type CreateProgramRequest struct {
	Title       string       `json:"title" validate:"required,max=255"`
	Description string       `json:"description" validate:"required,max=255"`
	Weeks       int          `json:"weeks" validate:"required"`
	Difficulty  domain.Level `json:"difficulty" validate:"required"`
	CategoryID  *uuid.UUID   `json:"category_id" validate:"omitempty"`
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

type ProgramTrainInput struct {
	TrainId    uuid.UUID `json:"train_id" validate:"required,uuid"`
	WeekNumber int       `json:"week_number" validate:"required"`
	DayOfWeek  int       `json:"day_of_week" validate:"required"`
	Position   int       `json:"position" validate:"required, minda=1"`
}
type AddProgramTrainsRequest struct {
	Trains []ProgramTrainInput `json:"trains" validate:"required"`
}

type ProgramTrainDB struct {
	domain.Train
	WeekNumber int `db:"week_number" json:"week_number"`
	DayOfWeek  int `db:"day_of_week" json:"day_of_week"`
	Position   int `db:"position" json:"position"`
}
