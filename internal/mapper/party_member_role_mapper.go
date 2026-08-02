package mapper

import (
	"github.com/reinp/event-platform/backend/internal/models/enum"
)

func UniquePartyRoles(
	roles []enum.PartyMemberRole,
) []enum.PartyMemberRole {

	uniqueRoles := make(
		[]enum.PartyMemberRole,
		0,
		len(roles),
	)

	seenRoles := make(
		map[enum.PartyMemberRole]struct{},
		len(roles),
	)

	for _, role := range roles {

		if _, exists := seenRoles[role]; exists {
			continue
		}

		seenRoles[role] = struct{}{}

		uniqueRoles = append(
			uniqueRoles,
			role,
		)
	}

	return uniqueRoles
}
