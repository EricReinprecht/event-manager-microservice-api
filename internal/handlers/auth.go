package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/responses"
	"github.com/reinp/event-platform/backend/internal/service"
)

func RequirePartyOwner(
	c *gin.Context,
	partyService *service.PartyService,
	partyID uuid.UUID,
) bool {

	userID, ok := getUserID(c)

	if !ok {
		responses.Unauthorized(c)
		return false
	}

	_, err := partyService.FindOwnedParty(
		c.Request.Context(),
		partyID,
		userID,
	)

	if err != nil {

		responses.Forbidden(c)

		return false
	}

	return true
}
