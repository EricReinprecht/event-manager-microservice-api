package models

import (
	"time"

	"github.com/google/uuid"
)

type PartyStage struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	PartyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_party_stage_name"`
	Party   Party     `gorm:"constraint:OnDelete:CASCADE"`

	Name string `gorm:"not null;uniqueIndex:idx_party_stage_name"`

	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
}
