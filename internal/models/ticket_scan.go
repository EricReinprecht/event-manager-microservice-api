package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"gorm.io/gorm"
)

type TicketScan struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	TicketID uuid.UUID
	Ticket   Ticket

	ScannedByID uuid.UUID
	ScannedBy   User

	ScannedAt time.Time

	Status enum.TicketScanStatus `gorm:"type:varchar(20);check:role IN ('PENDING','VERIFIED','REJECTED')"`

	VerifiedAt *time.Time

	VerifiedByID *uuid.UUID
	VerifiedBy   *User

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
