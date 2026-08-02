package party_member_repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type PartyMemberWriteRepository struct {
	transactionManager *database.TransactionManager
}

func NewPartyMemberWriteRepository(
	transactionManager *database.TransactionManager,
) *PartyMemberWriteRepository {

	return &PartyMemberWriteRepository{
		transactionManager: transactionManager,
	}
}

func (r *PartyMemberWriteRepository) SyncRoles(
	ctx context.Context,
	memberID uuid.UUID,
	roles []enum.PartyMemberRole,
) error {

	return r.transactionManager.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			if err := tx.
				Where(
					"party_member_id = ?",
					memberID,
				).
				Delete(
					&models.PartyMemberRole{},
				).
				Error(); err != nil {

				return err
			}

			if len(roles) == 0 {
				return nil
			}

			memberRoles := make(
				[]models.PartyMemberRole,
				0,
				len(roles),
			)

			for _, role := range roles {

				memberRoles = append(
					memberRoles,
					models.PartyMemberRole{
						ID:            uuid.New(),
						PartyMemberID: memberID,
						Role:          role,
					},
				)
			}

			return tx.
				Create(
					&memberRoles,
				).
				Error()
		},
	)
}
