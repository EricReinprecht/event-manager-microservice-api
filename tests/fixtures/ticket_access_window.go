package fixtures

import (
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
)

func TicketAccessWindow(ticketCategoryID uuid.UUID) models.TicketAccessWindow {
	now := time.Now().UTC()

	return models.TicketAccessWindow{
		ID:               uuid.New(),
		TicketCategoryID: ticketCategoryID,
		StartsAt:         now.Add(-time.Hour),
		EndsAt:           now.Add(time.Hour),
	}
}
