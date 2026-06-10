package domain

import (
	"github.com/google/uuid"
	"time"
)

type Level string

const (
	Easy   Level = "easy"
	Medium Level = "medium"
	Hard   Level = "hard"
)

type Program struct {
	ID         uuid.UUID  `json:"id"`
	Title      string     `json:"title"`
	Desc       string     `json:"-"`
	Weeks      int        `json:"weeks"`
	Difficult  string     `json:"difficulty"`
	IsPublic   bool       `json:"is_public"`
	CategoryID *uuid.UUID `json:"category_id"`
	ImageURL   string     `json:"image_url"`
	CreatedBy  uuid.UUID  `json:"-"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Version    int64      `json:"version"`
}
type ProgramCreatedEvent struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Weeks       int        `json:"weeks"`
	Difficulty  Level      `json:"difficulty"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
}
type ProgramAndTrainsCount struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Desc        string     `json:"-"`
	Weeks       int        `json:"weeks"`
	Difficult   string     `json:"difficulty"`
	IsPublic    bool       `json:"is_public"`
	CategoryID  *uuid.UUID `json:"category_id"`
	ImagePath   string     `json:"-"`
	ImageURL    string     `json:"image_url"`
	CreatedBy   uuid.UUID  `json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Version     int64      `json:"version"`
	TrainsCount int64      `json:"trains_count"`
}
