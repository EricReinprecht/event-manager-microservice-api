package fixtures

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
)

func PartyMember(
	userID uuid.UUID,
	partyID uuid.UUID,
) models.PartyMember {

	member := models.PartyMember{

		ID: uuid.New(),

		UserID: userID,

		PartyID: partyID,

		Roles: []models.PartyMemberRole{
			{
				ID: uuid.New(),

				Role: enum.RoleStaff,
			},
		},
	}

	return member
}
