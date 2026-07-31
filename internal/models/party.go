package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Party struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Title string `json:"title"`

	Description string `json:"description"`

	StartAt time.Time `json:"startAt"`

	EndAt time.Time `json:"endAt"`

	// Location
	LocationName string
	Street       string  `json:"street"`
	HouseNumber  string  `json:"houseNumber"`
	City         string  `json:"city"`
	Country      string  `json:"country"`
	PostalCode   string  `json:"postalCode"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Timezone     string  `json:"timezone"`

	ThumbnailID *uuid.UUID `json:"thumbnailId"`

	Thumbnail *Media `gorm:"foreignKey:ThumbnailID" json:"thumbnail"`

	Images []Media `gorm:"many2many:party_media;" json:"images"`

	TicketCategories []TicketCategory `json:"ticketCategories"`

	Categories []Category `gorm:"many2many:party_categories;" json:"categories"`

	OrganizerID uuid.UUID `json:"organizerId"`

	Organizer User `gorm:"foreignKey:OrganizerID" json:"organizer"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
