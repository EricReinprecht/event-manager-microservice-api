package models

import (
	"time"

	"github.com/google/uuid"
)

type RefundPolicy struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	TicketCategoryID uuid.UUID
	TicketCategory   TicketCategory

	// when this rule applies
	Until time.Time

	// percentage refunded
	Percentage int

	CreatedAt time.Time
}
