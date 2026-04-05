package constants

import "github.com/google/uuid"

type UserClaims struct {
	ID uuid.UUID `json:"id"`
}

type contextKey string

const UserClaimsKey contextKey = "userClaims"
