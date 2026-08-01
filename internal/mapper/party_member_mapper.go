package mapper

import (
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
)

func PartyMemberResponse(
	member models.PartyMember,
) dto.PartyMemberResponse {

	return dto.PartyMemberResponse{
		ID:        member.ID,
		UserID:    member.UserID,
		PartyID:   member.PartyID,
		Roles:     PartyMemberRoleResponses(member.Roles),
		CreatedAt: member.CreatedAt,
	}
}

func PartyMemberResponses(
	members []models.PartyMember,
) []dto.PartyMemberResponse {

	result := make(
		[]dto.PartyMemberResponse,
		0,
		len(members),
	)

	for _, member := range members {

		result = append(
			result,
			PartyMemberResponse(member),
		)
	}

	return result
}

func PartyMemberRoleResponses(
	roles []models.PartyMemberRole,
) []dto.PartyRoleResponse {

	result := make(
		[]dto.PartyRoleResponse,
		0,
		len(roles),
	)

	for _, role := range roles {

		result = append(
			result,
			dto.PartyRoleResponse{
				ID:   role.ID,
				Role: role.Role,
			},
		)
	}

	return result
}
