package models

import (
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	User   User      `gorm:"constraint:OnDelete:CASCADE"`

	StageName string

	Description string

	ImageID *uuid.UUID
	Image   *Media `gorm:"foreignKey:ImageID;constraint:OnDelete:SET NULL"`

	Verified bool `gorm:"not null;default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
