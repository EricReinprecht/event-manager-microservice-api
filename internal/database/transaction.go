package database

import "context"

type TransactionManager struct {
	db DBExecutor
}

func NewTransactionManager(
	db DBExecutor,
) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (t *TransactionManager) Transaction(
	ctx context.Context,
	fn func(tx DBExecutor) error,
) error {

	tx := t.db.
		WithContext(ctx).
		Begin()

	if err := tx.Error(); err != nil {
		return err
	}

	if err := fn(tx); err != nil {

		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}
