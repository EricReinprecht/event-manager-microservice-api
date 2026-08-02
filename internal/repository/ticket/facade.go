package ticket_repository

import "github.com/reinp/event-platform/backend/internal/database"

type Facade struct {
	Repository *TicketRepository
}

func NewFacade(db database.DBExecutor) *Facade {
	return &Facade{Repository: NewTicketRepository(db)}
}
