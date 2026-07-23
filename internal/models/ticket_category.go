package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketCategory struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Name string `gorm:"uniqueIndex:idx_party_ticket_category_name"`

	Price float64

	Amount int

	PartyID uuid.UUID `gorm:"uniqueIndex:idx_party_ticket_category_name"`

	Party Party

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
