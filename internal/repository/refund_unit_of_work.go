package repository

import (
	"context"

	"github.com/reinp/event-platform/backend/internal/database"

	purchaseRepository "github.com/reinp/event-platform/backend/internal/repository/purchase"
)

type RefundTransactionRepositories struct {
	Purchases *purchaseRepository.PurchaseRepository
	Tickets   *TicketRepository
}

type RefundUnitOfWork struct {
	transactionManager *database.TransactionManager
}

func NewRefundUnitOfWork(
	transactionManager *database.TransactionManager,
) *RefundUnitOfWork {

	return &RefundUnitOfWork{
		transactionManager: transactionManager,
	}
}

func (u *RefundUnitOfWork) Transaction(
	ctx context.Context,
	fn func(
		repositories *RefundTransactionRepositories,
	) error,
) error {

	return u.transactionManager.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			repositories :=
				&RefundTransactionRepositories{
					Purchases: purchaseRepository.NewPurchaseRepository(
						tx,
					),

					Tickets: NewTicketRepository(
						tx,
					),
				}

			return fn(
				repositories,
			)
		},
	)
}
