package domain

import (
	"github.com/google/uuid"
	"time"
)

type Train struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Type       string    `json:"type"`
	Duration   int64     `json:"duration"`
	IsPublic   bool      `json:"is_public"`
	Difficulty string    `json:"difficulty"`
	Calories   int64     `json:"calories"`
	CategoryId uuid.UUID `json:"category_id"`
	ImageURL   string    `json:"image_url"`
	CreatedBy  uuid.UUID `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type Exercise struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	VideoURL    string    `json:"video_url"`
}
type TrainExercise struct {
	ID          uuid.UUID `json:"exercise_id"`
	Title       string    `json:"exercise_title"`
	Description string    `json:"description"`
	VideoURL    string    `json:"video_url"`
	Difficulty  string    `json:"difficulty"`
	Steps       int       `json:"steps"`
	Sets        int       `json:"sets"`
	Position    int       `json:"position"`
}

type ProgramTrains struct {
	Train
	WeekNumber int `json:"week_number"`
	DayOfWeek  int `json:"day_of_week"`
	Position   int `json:"position"`
}

type ProgramAndTrain struct {
	Program ProgramTrains   `json:"program"`
	Trains  []ProgramTrains `json:"trains"`
}
