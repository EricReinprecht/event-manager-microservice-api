package models

import (
	"time"

	"github.com/google/uuid"
)

type StaffMember struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	PartyID uuid.UUID `gorm:"type:uuid;not null;index"`
	Party   Party     `gorm:"constraint:OnDelete:CASCADE"`

	UserID *uuid.UUID
	User   *User `gorm:"constraint:OnDelete:SET NULL"`

	FirstName string
	LastName  string

	Role        string
	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
}
