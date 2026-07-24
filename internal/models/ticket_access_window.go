package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketAccessWindow struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	TicketCategoryID uuid.UUID
	TicketCategory   TicketCategory

	StartsAt time.Time

	EndsAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
