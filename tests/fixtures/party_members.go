package fixtures

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
)

func PartyMember(
	userID uuid.UUID,
	partyID uuid.UUID,
	role enum.PartyRole,
) models.PartyMember {

	return models.PartyMember{

		ID: uuid.New(),

		UserID: userID,

		PartyID: partyID,

		Role: role,
	}
}
