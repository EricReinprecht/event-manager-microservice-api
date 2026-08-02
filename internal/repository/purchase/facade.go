package purchase_repository

import (
	"github.com/reinp/event-platform/backend/internal/database"
)

type Facade struct {
	Repository *PurchaseRepository
	Write      *PurchaseWriteRepository
}

func NewFacade(
	db database.DBExecutor,
	transactionManager *database.TransactionManager,
) *Facade {

	return &Facade{
		Repository: NewPurchaseRepository(
			db,
		),

		Write: NewPurchaseWriteRepository(
			transactionManager,
		),
	}
}
