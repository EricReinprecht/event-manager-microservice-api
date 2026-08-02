package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"gorm.io/gorm"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Email string `gorm:"uniqueIndex;not null"`

	Username string `gorm:"uniqueIndex"`

	PasswordHash string `gorm:"not null"`

	VerifiedAt *time.Time

	ProfileCompleted bool `gorm:"default:false"`

	AvatarID *uuid.UUID
	Avatar   *Media

	FirstName string
	LastName  string

	Roles []UserRole `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	ArtistProfile *Artist `gorm:"foreignKey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) HasRole(role enum.UserRole) bool {
	for _, assigned := range u.Roles {
		if assigned.Role == role {
			return true
		}
	}
	return false
}
