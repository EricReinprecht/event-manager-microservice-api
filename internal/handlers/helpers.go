package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/service"
)

func requirePartyOwner(
	c *gin.Context,
	partyService *service.PartyService,
	partyID uuid.UUID,
) bool {

	return RequirePartyOwner(
		c,
		partyService,
		partyID,
	)
}

func getUserID(
	c *gin.Context,
) (uuid.UUID, bool) {

	value, exists := c.Get("userID")

	if !exists {
		return uuid.Nil, false
	}

	switch id := value.(type) {

	case uuid.UUID:
		return id, true

	case string:

		parsed, err := uuid.Parse(id)

		if err != nil {
			return uuid.Nil, false
		}

		return parsed, true
	}

	return uuid.Nil, false
}
