package refresh_token_repository

import "github.com/reinp/event-platform/backend/internal/database"

type Facade struct {
	Repository *RefreshTokenRepository
}

func NewFacade(db database.DBExecutor) *Facade {
	return &Facade{Repository: NewRefreshTokenRepository(db)}
}
