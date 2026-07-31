package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateAccessWindowRequest struct {
	StartsAt time.Time `json:"startsAt"`

	EndsAt time.Time `json:"endsAt"`
}

type UpdateAccessWindowRequest struct {
	ID *uuid.UUID `json:"id"`

	StartsAt time.Time `json:"startsAt"`

	EndsAt time.Time `json:"endsAt"`
}

type AccessWindowResponse struct {
	ID uuid.UUID `json:"id"`

	StartsAt time.Time `json:"startsAt"`

	EndsAt time.Time `json:"endsAt"`
}
