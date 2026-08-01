package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateAccessWindowRequest struct {
	StartsAt time.Time `json:"startsAt" binding:"required"`

	EndsAt time.Time `json:"endsAt" binding:"required"`
}

type UpdateAccessWindowRequest struct {
	ID *uuid.UUID `json:"id"`

	StartsAt time.Time `json:"startsAt" binding:"omitempty"`

	EndsAt time.Time `json:"endsAt" binding:"omitempty"`
}

type AccessWindowResponse struct {
	ID uuid.UUID `json:"id"`

	StartsAt time.Time `json:"startsAt"`

	EndsAt time.Time `json:"endsAt"`
}
