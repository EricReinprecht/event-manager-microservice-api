package party_repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	baseRepository "github.com/reinp/event-platform/backend/internal/repository"
	ticketCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/ticketCategory"
)

type PartyWriteRepository struct {
	transactionManager *database.TransactionManager

	partyImages     *baseRepository.PartyImageRepository
	partyCategories *baseRepository.PartyCategoryRepository

	ticketCategories *ticketCategoryRepository.TicketCategoryWriteRepository
}

func NewPartyWriteRepository(
	transactionManager *database.TransactionManager,
	partyImages *baseRepository.PartyImageRepository,
	partyCategories *baseRepository.PartyCategoryRepository,
	ticketCategories *ticketCategoryRepository.TicketCategoryWriteRepository,
) *PartyWriteRepository {

	return &PartyWriteRepository{
		transactionManager: transactionManager,
		partyImages:        partyImages,
		partyCategories:    partyCategories,
		ticketCategories:   ticketCategories,
	}
}

func (r *PartyWriteRepository) CreateWithRelations(
	ctx context.Context,
	party *models.Party,
	categoryIDs []uuid.UUID,
	imageIDs []uuid.UUID,
) error {

	return r.transactionManager.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			parties := NewPartyRepository(tx)

			if err := parties.Create(
				ctx,
				party,
			); err != nil {
				return err
			}

			if err := r.partyCategories.Replace(
				tx,
				party.ID,
				categoryIDs,
			); err != nil {
				return err
			}

			if err := r.partyImages.Replace(
				tx,
				party.ID,
				imageIDs,
			); err != nil {
				return err
			}

			return createOrganizerMembership(
				tx,
				party.ID,
				party.OrganizerID,
			)
		},
	)
}

func (r *PartyWriteRepository) UpdateWithRelations(
	ctx context.Context,
	party *models.Party,
	categoryIDs []uuid.UUID,
	imageIDs []uuid.UUID,
	ticketCategories []models.TicketCategory,
) error {

	return r.transactionManager.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			parties := NewPartyRepository(tx)

			if err := parties.Update(
				ctx,
				party,
			); err != nil {
				return err
			}

			if err := r.partyCategories.Replace(
				tx,
				party.ID,
				categoryIDs,
			); err != nil {
				return err
			}

			if err := r.partyImages.Replace(
				tx,
				party.ID,
				imageIDs,
			); err != nil {
				return err
			}

			if err := r.ticketCategories.SyncCategories(
				tx,
				party.ID,
				ticketCategories,
			); err != nil {
				return err
			}

			return nil
		},
	)
}

func createOrganizerMembership(
	tx database.DBExecutor,
	partyID uuid.UUID,
	organizerID uuid.UUID,
) error {

	member := &models.PartyMember{
		ID:      uuid.New(),
		UserID:  organizerID,
		PartyID: partyID,
	}

	if err := tx.
		Create(member).
		Error(); err != nil {

		return err
	}

	role := &models.PartyMemberRole{
		ID:            uuid.New(),
		PartyMemberID: member.ID,
		Role:          enum.PartyRoleOrganizer,
	}

	return tx.
		Create(role).
		Error()
}

func (r *PartyWriteRepository) ReplaceImages(
	ctx context.Context,
	partyID uuid.UUID,
	imageIDs []uuid.UUID,
) error {

	return r.transactionManager.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			return r.partyImages.Replace(
				tx,
				partyID,
				imageIDs,
			)
		},
	)
}

func (r *PartyWriteRepository) ReplaceCategories(
	ctx context.Context,
	partyID uuid.UUID,
	categoryIDs []uuid.UUID,
) error {

	return r.transactionManager.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			return r.partyCategories.Replace(
				tx,
				partyID,
				categoryIDs,
			)
		},
	)
}
