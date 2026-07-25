package models

import (
	"time"

	"github.com/google/uuid"
)

type PartyMember struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"uniqueIndex:idx_party_member_user_party"`
	User   User

	PartyID uuid.UUID `gorm:"uniqueIndex:idx_party_member_user_party"`
	Party   Party

	Roles []PartyMemberRole `gorm:"foreignKey:PartyMemberID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
