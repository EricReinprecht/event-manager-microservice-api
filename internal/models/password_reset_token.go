package models

import (
	"time"

	"github.com/google/uuid"
)

type PasswordResetToken struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID

	User User `gorm:"foreignKey:UserID"`

	TokenHash string `gorm:"uniqueIndex;not null"`

	ExpiresAt time.Time

	InvalidatedAt *time.Time

	UsedAt *time.Time

	CreatedAt time.Time
}
