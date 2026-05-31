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
}
