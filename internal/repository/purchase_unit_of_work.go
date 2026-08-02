package repository

import (
	"context"

	"github.com/reinp/event-platform/backend/internal/database"
)

type PurchaseTransactionRepositories struct {
	Purchases        *PurchaseRepository
	TicketCategories *TicketCategoryRepository
	Tickets          *TicketRepository
}

type PurchaseUnitOfWork struct {
	transactionManager *database.TransactionManager
}

func NewPurchaseUnitOfWork(
	transactionManager *database.TransactionManager,
) *PurchaseUnitOfWork {

	return &PurchaseUnitOfWork{
		transactionManager: transactionManager,
	}
}

func (u *PurchaseUnitOfWork) Transaction(
	ctx context.Context,
	fn func(
		repositories *PurchaseTransactionRepositories,
	) error,
) error {

	return u.transactionManager.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			repositories := &PurchaseTransactionRepositories{
				Purchases: NewPurchaseRepository(
					tx,
				),

				TicketCategories: NewTicketCategoryRepository(
					tx,
				),

				Tickets: NewTicketRepository(
					tx,
				),
			}

			return fn(repositories)
		},
	)
}
