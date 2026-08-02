package password_reset_token_repository

import "github.com/reinp/event-platform/backend/internal/database"

type Facade struct {
	Repository *PasswordResetTokenRepository
}

func NewFacade(db database.DBExecutor) *Facade {
	return &Facade{Repository: NewPasswordResetTokenRepository(db)}
}
