package helpers

import (
	appModels "github.com/reinp/event-platform/backend/internal/models"
)

type VerificationScenario struct {
	Staff appModels.User

	Customer appModels.User

	Party appModels.Party

	TicketCategory appModels.TicketCategory

	Ticket appModels.Ticket

	Window appModels.TicketAccessWindow

	Scan appModels.TicketScan
}
