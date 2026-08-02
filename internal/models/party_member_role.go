package models

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type PartyMemberRole struct {
	ID uuid.UUID

	PartyMemberID uuid.UUID
	PartyMember   PartyMember

	Role enum.PartyMemberRole
}
