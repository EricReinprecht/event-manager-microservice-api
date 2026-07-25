package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type Purchase struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`
	User   User

	PartyID uuid.UUID `gorm:"type:uuid;not null"`
	Party   Party

	Status enum.PurchaseStatus `gorm:"type:varchar(20);check:status IN ('PENDING','PAID','FAILED','CANCELED','REFUNDED')"`

	PaymentProvider string
	PaymentID       string

	ExpiresAt time.Time

	TotalPrice int64

	Items []PurchaseItem `gorm:"foreignKey:PurchaseID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
