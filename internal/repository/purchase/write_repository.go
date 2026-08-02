package purchase_repository

import (
	"github.com/reinp/event-platform/backend/internal/database"
)

type PurchaseWriteRepository struct {
	transactionManager *database.TransactionManager
}

func NewPurchaseWriteRepository(
	transactionManager *database.TransactionManager,
) *PurchaseWriteRepository {

	return &PurchaseWriteRepository{
		transactionManager: transactionManager,
	}
}
