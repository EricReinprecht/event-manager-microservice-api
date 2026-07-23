package models

import "github.com/google/uuid"

type PartyMedia struct {
	PartyID uuid.UUID `gorm:"type:uuid;primaryKey"`

	MediaID uuid.UUID `gorm:"type:uuid;primaryKey"`

	IsThumbnail bool `gorm:"default:false"`

	Position int

	Party Party

	Media Media
}
