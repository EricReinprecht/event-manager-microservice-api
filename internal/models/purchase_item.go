package models

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseItem struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	PurchaseID uuid.UUID

	Purchase Purchase

	TicketCategoryID uuid.UUID

	TicketCategory TicketCategory

	Quantity int

	CreatedAt time.Time
}
