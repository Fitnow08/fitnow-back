package train

import "github.com/go-playground/validator/v10"

type CreateTrainRequest struct {
	Title      string `json:"title" validate:"required"`
	Type       string `json:"type" validate:"required,oneof=strength cardio stretching"`
	Duration   int64  `json:"duration" validate:"required,min=1"`
	IsPublic   bool   `json:"is_public"`
	Difficulty string `json:"difficulty" validate:"required,oneof=easy medium hard"`
	Calories   int64  `json:"calories" validate:"required,min=0"`
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

func NewValidator() *validator.Validate {
	return validator.New()
}
