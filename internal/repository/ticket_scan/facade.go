package ticket_scan_repository

import "github.com/reinp/event-platform/backend/internal/database"

type Facade struct {
	Repository *TicketScanRepository
}

func NewFacade(db database.DBExecutor) *Facade {
	return &Facade{Repository: NewTicketScanRepository(db)}
}
