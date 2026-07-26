package fixtures

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
)

func Purchase(
	userID uuid.UUID,
	partyID uuid.UUID,
) models.Purchase {

	return models.Purchase{

		ID: uuid.New(),

		UserID: userID,

		PartyID: partyID,

		Status: enum.PurchaseStatusPending,

		TotalPrice: 1000,
	}
}
