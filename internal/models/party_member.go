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

	Roles []PartyMemberRole `gorm:"foreignKey:PartyMemberID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *PartyMember) HasRole(
	role enum.PartyMemberRole,
) bool {

	for _, r := range m.Roles {

		if r.Role == role {
			return true
		}
	}

	return false
}
