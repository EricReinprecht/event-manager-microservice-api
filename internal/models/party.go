package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Party struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Title string

	Description string

	StartAt time.Time
	EndAt   time.Time

	Location string

	ThumbnailID *uuid.UUID

	Thumbnail *Media `gorm:"foreignKey:ThumbnailID"`

	Images []Media `gorm:"many2many:party_media;"`

	TicketCategories []TicketCategory

	Categories []Category `gorm:"many2many:party_categories;"`

	OrganizerID uuid.UUID

	Organizer User `gorm:"foreignKey:OrganizerID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
