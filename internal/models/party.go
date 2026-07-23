package models

import (
	"time"

	"github.com/google/uuid"
)

type Party struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Title string `gorm:"not null"`

	Description string

	Date time.Time `gorm:"not null"`

	Location string

	OrganizerID uuid.UUID `gorm:"type:uuid;not null"`

	Organizer User `gorm:"foreignKey:OrganizerID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
