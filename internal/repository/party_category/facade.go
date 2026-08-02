package party_category_repository

import (
	"github.com/reinp/event-platform/backend/internal/database"
)

type Facade struct {
	Repository *PartyCategoryRepository
	Write      *PartyCategoryWriteRepository
}

func NewFacade(
	db database.DBExecutor,
) *Facade {

	return &Facade{
		Repository: NewPartyCategoryRepository(
			db,
		),

		Write: NewPartyCategoryWriteRepository(),
	}
}
