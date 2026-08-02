package models

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type PartyMemberRole struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	PartyMemberID uuid.UUID `gorm:"type:uuid;not null"`

	Role enum.PartyMemberRole `gorm:"type:varchar(20);not null"`
}
