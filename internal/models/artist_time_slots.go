package models

import (
	"time"

	"github.com/google/uuid"
)

type ArtistSlot struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	PartyID uuid.UUID `gorm:"type:uuid;not null;index"`
	Party   Party     `gorm:"constraint:OnDelete:CASCADE"`

	ArtistID uuid.UUID `gorm:"type:uuid;not null;index"`
	Artist   Artist    `gorm:"constraint:OnDelete:CASCADE"`

	StageID uuid.UUID  `gorm:"type:uuid;not null;index"`
	Stage   PartyStage `gorm:"constraint:OnDelete:CASCADE"`

	StartAt time.Time

	EndAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
