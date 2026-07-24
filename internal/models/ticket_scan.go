package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketScan struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	TicketID uuid.UUID
	Ticket   Ticket

	ScannedByID uuid.UUID
	ScannedBy   User

	ScannedAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
