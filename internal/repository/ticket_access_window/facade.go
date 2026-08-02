package ticket_access_window_repository

import "github.com/reinp/event-platform/backend/internal/database"

type Facade struct {
	Repository *TicketAccessWindowRepository
}

func NewFacade(db database.DBExecutor) *Facade {
	return &Facade{Repository: NewTicketAccessWindowRepository(db)}
}
