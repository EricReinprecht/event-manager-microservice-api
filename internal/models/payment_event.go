package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentEvent struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Provider string `gorm:"not null"`

	EventID string `gorm:"uniqueIndex;not null"`

	Type string `gorm:"not null"`

	Payload string `gorm:"type:text;not null"`

	Processed bool `gorm:"not null;default:false"`

	CreatedAt time.Time

	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
