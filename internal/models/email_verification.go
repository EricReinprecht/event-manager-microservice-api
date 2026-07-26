package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`

	Token string `gorm:"uniqueIndex;not null"`

	ExpiresAt time.Time

	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID"`
}
