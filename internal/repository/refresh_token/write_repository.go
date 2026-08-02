package refresh_token_repository

import (
	"context"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type RefreshTokenWriteRepository struct {
	transactionManager *database.TransactionManager
}

func NewRefreshTokenWriteRepository(transactionManager *database.TransactionManager) *RefreshTokenWriteRepository {
	return &RefreshTokenWriteRepository{transactionManager: transactionManager}
}

func (r *RefreshTokenWriteRepository) Rotate(
	ctx context.Context,
	current *models.RefreshToken,
	next *models.RefreshToken,
) error {
	return r.transactionManager.Transaction(ctx, func(tx database.DBExecutor) error {
		tokens := NewRefreshTokenRepository(tx)
		if err := tokens.Revoke(current.ID); err != nil {
			return err
		}
		return tokens.Create(next)
	})
}
