package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID
	User   User

	Token string `gorm:"uniqueIndex;not null"`

	ExpiresAt time.Time

	VerifiedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
