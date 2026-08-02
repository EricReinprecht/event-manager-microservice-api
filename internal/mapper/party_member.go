package mapper

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
)

func CreatePartyMember(
	partyID uuid.UUID,
	req dto.CreatePartyMemberRequest,
) *models.PartyMember {

	member := &models.PartyMember{
		ID:      uuid.New(),
		UserID:  req.UserID,
		PartyID: partyID,
		Roles: make(
			[]models.PartyMemberRole,
			0,
			len(req.Roles),
		),
	}

	for _, role := range req.Roles {

		member.Roles = append(
			member.Roles,
			models.PartyMemberRole{
				ID:            uuid.New(),
				PartyMemberID: member.ID,
				Role:          role,
			},
		)
	}

	return member
}
