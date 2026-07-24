package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Code string `gorm:"uniqueIndex;not null"`

	TicketCategoryID uuid.UUID
	TicketCategory   TicketCategory

	UserID uuid.UUID
	User   User

	Scans []TicketScan

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
