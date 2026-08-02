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

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`

	Roles enum.UserRole `gorm:"type:varchar(20);not null;default:'DEFAULT';check:ticket_status_check,status IN ('PREMIUM','PLATIUM')"`
}
