package repository

import (
	"context"

	"github.com/reinp/event-platform/backend/internal/database"
)

type TicketTransactionRepositories struct {
	Tickets       *TicketRepository
	Scans         *TicketScanRepository
	AccessWindows *TicketAccessWindowRepository
}

type TicketUnitOfWork struct {
	transactionManager *database.TransactionManager
}

func NewTicketUnitOfWork(
	transactionManager *database.TransactionManager,
) *TicketUnitOfWork {

	return &TicketUnitOfWork{
		transactionManager: transactionManager,
	}
}

func (u *TicketUnitOfWork) Transaction(
	ctx context.Context,
	fn func(repositories *TicketTransactionRepositories) error,
) error {

	return u.transactionManager.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			repositories := &TicketTransactionRepositories{
				Tickets: NewTicketRepository(
					tx,
				),

				Scans: NewTicketScanRepository(
					tx,
				),

				AccessWindows: NewTicketAccessWindowRepository(
					tx,
				),
			}

			return fn(repositories)
		},
	)
}
