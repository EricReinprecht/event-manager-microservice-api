package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type PartyMember struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"uniqueIndex:idx_party_member_user_party"`
	User   User

	PartyID uuid.UUID `gorm:"uniqueIndex:idx_party_member_user_party"`
	Party   Party

	Role enum.PartyRole `gorm:"type:varchar(20);check:role IN ('ORGANIZER','ADMIN','STAFF','ATTENDEE')"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
