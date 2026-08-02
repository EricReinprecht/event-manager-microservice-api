package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type PartyWriteRepository struct {
	transactionManager *database.TransactionManager
	partyImages        *PartyImageRepository
	partyCategories    *PartyCategoryRepository
	ticketCategories   *TicketCategoryWriteRepository
}

func NewPartyWriteRepository(
	transactionManager *database.TransactionManager,
	partyImages *PartyImageRepository,
	partyCategories *PartyCategoryRepository,
	ticketCategories *TicketCategoryWriteRepository,
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

			if err := tx.
				Create(
					party,
				).
				Error(); err != nil {

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

			if err := tx.
				Model(
					&models.Party{},
				).
				Where(
					"id = ?",
					party.ID,
				).
				Updates(
					partyUpdateValues(
						party,
					),
				).
				Error(); err != nil {

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
		Create(
			member,
		).
		Error(); err != nil {

		return err
	}

	role := &models.PartyMemberRole{
		ID:            uuid.New(),
		PartyMemberID: member.ID,
		Role:          enum.PartyRoleOrganizer,
	}

	return tx.
		Create(
			role,
		).
		Error()
}

func partyUpdateValues(
	party *models.Party,
) map[string]any {

	return map[string]any{
		"title": party.Title,

		"description": party.Description,

		"thumbnail_id": party.ThumbnailID,

		"location_name": party.LocationName,

		"street": party.Street,

		"house_number": party.HouseNumber,

		"city": party.City,

		"country": party.Country,

		"postal_code": party.PostalCode,

		"latitude": party.Latitude,

		"longitude": party.Longitude,

		"timezone": party.Timezone,

		"start_at": party.StartAt,

		"end_at": party.EndAt,
	}
}
