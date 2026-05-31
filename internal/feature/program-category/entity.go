package programcategory

import (
	"github.com/google/uuid"
	"time"
)

type ProgramCategoryDB struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"title"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
	Version   int       `db:"version"`
}

type GetAllProgramCategoryResponse struct {
	Categories []ProgramCategoryDB `json:"categories"`
}

type CreateProgramCategoryRequest struct {
	Title string `json:"title" validate:"required"`
}

type CreateProgramCategoryResponse struct {
	Category *ProgramCategoryDB `json:"category"`
}

type UpdateProgramCategoryRequest struct {
	Title string `json:"title" validate:"required"`
}

type UpdateProgramCategoryResponse struct {
	Category *ProgramCategoryDB `json:"category"`
}
