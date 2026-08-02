package media_repository

import "github.com/reinp/event-platform/backend/internal/database"

type Facade struct {
	Repository *MediaRepository
}

func NewFacade(db database.DBExecutor) *Facade {
	return &Facade{Repository: NewMediaRepository(db)}
}
