package refresh_token_repository

import "github.com/reinp/event-platform/backend/internal/database"

type Facade struct {
	Repository *RefreshTokenRepository
	Write      *RefreshTokenWriteRepository
}

func NewFacade(db database.DBExecutor, transactionManager *database.TransactionManager) *Facade {
	return &Facade{
		Repository: NewRefreshTokenRepository(db),
		Write:      NewRefreshTokenWriteRepository(transactionManager),
	}
}
