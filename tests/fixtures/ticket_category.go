package fixtures

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
)

func TicketCategory(
	partyID uuid.UUID,
) models.TicketCategory {

	return models.TicketCategory{
		ID: uuid.New(),

		PartyID: partyID,

		Name: "General Admission",

		Price: 1000,

		Capacity: 100,
	}
}
