package payment_event_repository

import "github.com/reinp/event-platform/backend/internal/database"

type Facade struct {
	Repository *PaymentEventRepository
}

func NewFacade(db database.DBExecutor) *Facade {
	return &Facade{Repository: NewPaymentEventRepository(db)}
}
