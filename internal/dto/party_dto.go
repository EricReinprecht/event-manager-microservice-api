package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePartyRequest struct {
	Title string `json:"title"`

	Description string `json:"description"`

	Location string `json:"location"`

	StartAt time.Time `json:"start_at" binding:"required"`

	EndAt time.Time `json:"end_at" binding:"required"`

	CategoryID uuid.UUID `json:"category_id" binding:"required"`

	ThumbnailID *uuid.UUID `json:"thumbnail_id"`

	ImageIDs []uuid.UUID `json:"image_ids"`
}

type UpdatePartyRequest struct {
	Title string `json:"title"`

	Description string `json:"description"`

	Location string `json:"location"`

	StartAt time.Time `json:"start_at" binding:"required"`

	EndAt time.Time `json:"end_at" binding:"required"`

	CategoryID uuid.UUID `json:"category_id" binding:"required"`

	ThumbnailID *uuid.UUID `json:"thumbnail_id"`

	ImageIDs []uuid.UUID `json:"image_ids"`
}
