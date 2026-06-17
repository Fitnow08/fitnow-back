package traincategory

import (
	"github.com/google/uuid"
	"time"
)

type TrainCategoryDB struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"title" json:"title"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Version   int       `db:"version" json:"version"`
}

type GetAllTrainCategoryResponse struct {
	Categories []TrainCategoryDB `json:"categories"`
}
type CreateTrainCategoryRequest struct {
	Title string `json:"title" validate:"required"`
}

type CreateTrainCategoryResponse struct {
	Category *TrainCategoryDB `json:"category"`
}

type UpdateTrainCategoryRequest struct {
	Title string `json:"title" validate:"required"`
}
type UpdateTrainCategoryResponse struct {
	Category *TrainCategoryDB `json:"category"`
}
