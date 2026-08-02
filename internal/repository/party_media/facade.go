package party_media_repository

import "github.com/reinp/event-platform/backend/internal/database"

type Facade struct {
	Repository *PartyMediaRepository
}

func NewFacade(db database.DBExecutor) *Facade {
	return &Facade{Repository: NewPartyMediaRepository(db)}
}
