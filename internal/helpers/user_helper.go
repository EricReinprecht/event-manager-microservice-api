package helpers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
)

type PartyOwnerChecker interface {
	FindOwnedParty(
		ctx context.Context,
		partyID uuid.UUID,
		userID uuid.UUID,
	) (*models.Party, error)
}

func RequireUserID(
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

func RequirePartyOwner(
	c *gin.Context,
	partyService PartyOwnerChecker,
	partyID uuid.UUID,
) bool {

	userID, ok := RequireUserID(c)

	if !ok {

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid user",
			},
		)

		return false
	}

	_, err := partyService.FindOwnedParty(
		c.Request.Context(),
		partyID,
		userID,
	)

	if err != nil {

		c.JSON(
			http.StatusForbidden,
			gin.H{
				"error": "not allowed",
			},
		)

		return false
	}

	return true
}
