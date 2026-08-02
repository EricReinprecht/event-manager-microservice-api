package user_repository

import "github.com/reinp/event-platform/backend/internal/database"

type Facade struct {
	Repository *UserRepository
}

func NewFacade(db database.DBExecutor) *Facade {
	return &Facade{Repository: NewUserRepository(db)}
}
