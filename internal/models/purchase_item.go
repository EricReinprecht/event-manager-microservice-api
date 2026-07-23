package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseItem struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	PurchaseID uuid.UUID
	Purchase   Purchase

	TicketCategoryID uuid.UUID
	TicketCategory   TicketCategory

	Quantity int

	Price int64

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
