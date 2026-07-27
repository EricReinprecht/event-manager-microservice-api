package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Email string `gorm:"uniqueIndex;not null"`

	VerifiedAt *time.Time

	Username string `gorm:"uniqueIndex"`

	AvatarID *uuid.UUID
	Avatar   *Media

	PasswordHash string `gorm:"not null"`

	FirstName string
	LastName  string

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
