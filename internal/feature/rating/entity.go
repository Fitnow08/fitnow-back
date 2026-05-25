package rating

import (
	"github.com/google/uuid"
	"time"
)

type RatingDB struct {
	UserId    uuid.UUID `db:"userId"`
	TrainId   uuid.UUID `db:"trainId"`
	Rating    int       `db:"rating"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateRatingRequest struct {
	Rating int `json:"rating"`
}
