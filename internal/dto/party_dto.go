package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePartyRequest struct {
	Title string `json:"title" binding:"required,min=1,max=100"`

	Description string `json:"description" binding:"max=2000"`

	LocationName string `json:"locationName" binding:"required,max=150"`

	Location PartyLocation `json:"location" binding:"required"`

	StartAt time.Time `json:"startAt" binding:"required"`

	EndAt time.Time `json:"endAt" binding:"required"`

	CategoryIDs []string `json:"categories"`

	ThumbnailID *uuid.UUID `json:"thumbnailId"`

	ImageIDs []uuid.UUID `json:"imageIds"`

	TicketCategories []CreateTicketCategoryRequest `json:"ticketCategories" binding:"dive"`
}

type UpdatePartyRequest struct {
	Title string `json:"title" binding:"required,min=1,max=100"`

	Description string `json:"description" binding:"max=2000"`

	LocationName string `json:"locationName" binding:"required,max=150"`

	Location PartyLocation `json:"location" binding:"required"`

	StartAt time.Time `json:"startAt" binding:"required"`

	EndAt time.Time `json:"endAt" binding:"required"`

	CategoryIDs []string `json:"categories"`

	ThumbnailID *uuid.UUID `json:"thumbnailId"`

	ImageIDs []uuid.UUID `json:"imageIds"`

	TicketCategories []UpdateTicketCategoryRequest `json:"ticketCategories" binding:"dive"`
}

type PartyResponse struct {
	ID uuid.UUID `json:"id"`

	Title string `json:"title"`

	Description string `json:"description"`

	LocationName string `json:"locationName"`

	Location PartyLocation `json:"location"`

	StartAt time.Time `json:"startAt"`

	EndAt time.Time `json:"endAt"`

	ThumbnailID *uuid.UUID `json:"thumbnailId"`

	OrganizerID uuid.UUID `json:"organizerId"`

	Categories []CategoryResponse `json:"categories"`

	TicketCategories []TicketCategoryResponse `json:"ticketCategories"`
}
