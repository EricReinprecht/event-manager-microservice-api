package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseItem struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	PurchaseID uuid.UUID `gorm:"type:uuid;not null"`
	Purchase   Purchase  `gorm:"foreignKey:PurchaseID"`

	TicketCategoryID uuid.UUID      `gorm:"type:uuid;not null"`
	TicketCategory   TicketCategory `gorm:"foreignKey:TicketCategoryID"`

	Quantity int `gorm:"not null"`

	UnitPrice int64 `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
