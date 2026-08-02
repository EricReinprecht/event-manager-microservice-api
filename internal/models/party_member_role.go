package models

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type PartyMemberRole struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	PartyMemberID uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex:idx_party_member_role"`
	PartyMember   PartyMember `gorm:"constraint:OnDelete:CASCADE"`

	Role enum.PartyMemberRole `gorm:"type:varchar(20);not null;uniqueIndex:idx_party_member_role;check:party_member_role_check,role IN ('ORGANIZER','ADMIN','REFUNDER','SCANNER')"`
}
