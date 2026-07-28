package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`

	TokenHash string `gorm:"not null;uniqueIndex"`

	FamilyID uuid.UUID `gorm:"type:uuid;not null;index"`

	UserAgent string

	IPAddress string

	DeviceName string

	ExpiresAt time.Time `gorm:"not null"`

	RevokedAt *time.Time

	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID"`
}
