package fixtures

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
)

func Ticket() models.Ticket {
	return models.Ticket{
		ID:     uuid.New(),
		Code:   uuid.NewString(),
		Status: enum.TicketStatusActive,
	}
}
