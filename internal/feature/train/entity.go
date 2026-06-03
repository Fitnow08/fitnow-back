package train

import (
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type TrainDB struct {
	ID         uuid.UUID `db:"id"`
	Title      string    `db:"title"`
	Type       string    `db:"type"`
	Duration   int64     `db:"duration"`
	IsPublic   bool      `db:"is_public"`
	Difficulty string    `db:"difficulty"`
	CreatedBy  uuid.UUID `db:"created_by"`
	CreatedAt  time.Time `db:"created_at"`
	CategoryId uuid.UUID `db:"category_id"`
	Calories   int64     `db:"calories"`
	ImagePath  string    `db:"image_path"`
	UpdatedAt  time.Time `db:"updated_at"`
	Version    int64     `db:"version"`
}
type CreateTrainRequest struct {
	Title      string    `json:"title" validate:"required"`
	Type       string    `json:"type" validate:"required,oneof=strength cardio stretching"`
	Duration   int64     `json:"duration" validate:"required,min=1"`
	IsPublic   bool      `json:"is_public"`
	Difficulty string    `json:"difficulty" validate:"required,oneof=easy medium hard"`
	Calories   int64     `json:"calories" validate:"required,min=0"`
	CategoryId uuid.UUID `json:"category_id" validate:"required,uuid"`
}

type UpdateTrainRequest struct {
	Title      string `json:"title"`
	Type       string `json:"type" validate:"omitempty,oneof=strength cardio stretching"`
	Duration   int64  `json:"duration" validate:"omitempty,min=1"`
	IsPublic   *bool  `json:"is_public"`
	Difficulty string `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	Calories   int64  `json:"calories" validate:"omitempty,min=0"`
}

type CreateExerciseRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}
type AllTrainsParams struct {
	Page       uint64    `json:"page"`
	Limit      uint64    `json:"size"`
	CategoryId uuid.UUID `json:"category_id"`
	Text       string    `json:"text"`
}
type TrainExerciseInput struct {
	ExerciseID uuid.UUID `json:"exercise_id" validate:"required"`
	Steps      int       `json:"steps" validate:"required,min=1"`
	Sets       int       `json:"sets" validate:"required,min=1"`
	Position   int       `json:"position" validate:"required"`
}
type AddTrainExerciseRequest struct {
	Exercises []TrainExerciseInput `json:"exercises" validate:"required,min=1,dive"`
}

type TrainAndExercises struct {
	Train     domain.Train           `json:"train"`
	Exercises []domain.TrainExercise `json:"exercises"`
}

func NewValidator() *validator.Validate {
	return validator.New()
}
