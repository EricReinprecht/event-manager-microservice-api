package models

import (
	"time"

	"github.com/google/uuid"
)

type RefundPolicy struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	TicketCategoryID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	Until time.Time

	Percentage int

	CreatedAt time.Time
}
