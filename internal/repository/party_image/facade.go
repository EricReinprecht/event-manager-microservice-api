package party_image_repository

import "github.com/reinp/event-platform/backend/internal/database"

type Facade struct {
	Repository *PartyImageRepository
}

func NewFacade(db database.DBExecutor) *Facade {
	return &Facade{Repository: NewPartyImageRepository(db)}
}
