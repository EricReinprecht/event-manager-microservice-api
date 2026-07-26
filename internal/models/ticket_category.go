package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketCategory struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Name string

	Price int64

	Capacity int

	PartyID uuid.UUID

	Party Party

	RequiresVerification bool

	AccessWindows []TicketAccessWindow

	RefundRequiresApproval bool

	RefundPolicyID *uuid.UUID
	RefundPolicy   *RefundPolicy

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
