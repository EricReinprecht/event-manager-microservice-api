package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type Ticket struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Code string `gorm:"uniqueIndex;not null"`

	Status enum.TicketStatus `gorm:"type:varchar(20);not null;default:'ACTIVE';check:ticket_status_check,status IN ('ACTIVE','CANCELLED')"`

	TicketCategoryID uuid.UUID
	TicketCategory   TicketCategory

	UserID uuid.UUID
	User   User

	Scans []TicketScan

	PurchaseID uuid.UUID

	Purchase Purchase

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
