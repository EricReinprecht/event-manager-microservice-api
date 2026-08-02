package ticket_category_repository

import (
	"github.com/reinp/event-platform/backend/internal/database"
)

type Facade struct {
	Repository *TicketCategoryRepository
	Write      *TicketCategoryWriteRepository
}

func NewFacade(
	db database.DBExecutor,
) *Facade {

	return &Facade{
		Repository: NewTicketCategoryRepository(
			db,
		),

		Write: NewTicketCategoryWriteRepository(),
	}
}
