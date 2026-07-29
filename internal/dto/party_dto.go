package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePartyRequest struct {
	Title string `json:"title"`

	Description string `json:"description"`

	Location string `json:"location"`

	StartAt time.Time `json:"startAt" binding:"required"`

	EndAt time.Time `json:"endAt" binding:"required"`

	CategoryIDs []uuid.UUID `json:"categoryIds" binding:"required,min=1"`

	ThumbnailID *uuid.UUID `json:"thumbnailId"`

	ImageIDs []uuid.UUID `json:"imageIds"`
}

type UpdatePartyRequest struct {
	Title string `json:"title"`

	Description string `json:"description"`

	Location string `json:"location"`

	StartAt time.Time `json:"startAt" binding:"required"`

	EndAt time.Time `json:"endAt" binding:"required"`

	CategoryIDs []uuid.UUID `json:"categoryIds" binding:"required,min=1"`

	ThumbnailID *uuid.UUID `json:"thumbnailId"`

	ImageIDs []uuid.UUID `json:"imageIds"`
}
