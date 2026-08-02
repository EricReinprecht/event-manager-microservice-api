package email_verification_repository

import "github.com/reinp/event-platform/backend/internal/database"

type Facade struct {
	Repository *EmailVerificationRepository
}

func NewFacade(db database.DBExecutor) *Facade {
	return &Facade{Repository: NewEmailVerificationRepository(db)}
}
