package service

import (
	"context"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PartyCRUDService struct {
	parties          *repository.PartyRepository
	images           *repository.PartyImageRepository
	members          *repository.PartyMemberRepository
	roles            *repository.PartyMemberRoleRepository
	categories       *repository.CategoryRepository
	partyCategories  *repository.PartyCategoryRepository
	media            *repository.MediaRepository
	ticketCategories *repository.TicketCategoryRepository

	tx *database.TransactionManager
}

func NewPartyCRUDService(
	parties *repository.PartyRepository,
	images *repository.PartyImageRepository,
	members *repository.PartyMemberRepository,
	roles *repository.PartyMemberRoleRepository,
	categories *repository.CategoryRepository,
	partyCategories *repository.PartyCategoryRepository,
	media *repository.MediaRepository,
	ticketCategories *repository.TicketCategoryRepository,
	tx *database.TransactionManager,
) *PartyCRUDService {

	return &PartyCRUDService{
		parties:          parties,
		images:           images,
		members:          members,
		roles:            roles,
		categories:       categories,
		partyCategories:  partyCategories,
		media:            media,
		ticketCategories: ticketCategories,
		tx:               tx,
	}
}

func (s *PartyCRUDService) Create(
	ctx context.Context,
	party *models.Party,
	imageIDs []uuid.UUID,
) error {

	for _, category := range party.Categories {

		_, err := s.categories.FindByID(
			ctx,
			category.ID,
		)

		if err != nil {

			return appErrors.ErrCategoryNotFound
		}
	}

	if party.ThumbnailID != nil {

		_, err := s.media.FindByID(
			ctx,
			*party.ThumbnailID,
		)

		if err != nil {

			return appErrors.ErrMediaNotFound
		}
	}

	for _, imageID := range imageIDs {

		_, err := s.media.FindByID(
			ctx,
			imageID,
		)

		if err != nil {

			return appErrors.ErrMediaNotFound
		}
	}

	return s.tx.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			if err := s.parties.Create(
				tx,
				party,
			); err != nil {

				return err
			}

			if err := tx.
				Model(party).
				Association("Categories").
				Replace(
					party.Categories,
				); err != nil {

				return err
			}

			if err := s.images.Replace(
				tx,
				party.ID,
				imageIDs,
			); err != nil {

				return err
			}

			member := &models.PartyMember{
				ID: uuid.New(),

				UserID: party.OrganizerID,

				PartyID: party.ID,
			}

			if err := s.members.Create(
				tx,
				member,
			); err != nil {

				return err
			}

			role := &models.PartyMemberRole{
				ID: uuid.New(),

				PartyMemberID: member.ID,

				Role: enum.RoleOrganizer,
			}

			return s.roles.Create(
				tx,
				role,
			)
		},
	)
}

func (s *PartyCRUDService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Party, error) {

	return s.parties.FindByID(
		ctx,
		id,
	)
}

func (s *PartyCRUDService) Update(
	ctx context.Context,
	party *models.Party,
) error {

	return s.parties.Update(
		ctx,
		party,
	)
}

func (s *PartyCRUDService) UpdateImages(
	ctx context.Context,
	partyID uuid.UUID,
	imageIDs []uuid.UUID,
) error {

	return s.tx.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			return s.images.Replace(
				tx,
				partyID,
				imageIDs,
			)
		},
	)
}

func (s *PartyCRUDService) UpdateCategories(
	ctx context.Context,
	partyID uuid.UUID,
	categoryIDs []uuid.UUID,
) error {

	return s.tx.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			return s.partyCategories.Replace(
				tx,
				partyID,
				categoryIDs,
			)
		},
	)
}

func (s *PartyCRUDService) Delete(
	ctx context.Context,
	party *models.Party,
) error {

	return s.parties.Delete(
		ctx,
		party,
	)
}

func (s *PartyCRUDService) UpdateRelations(
	ctx context.Context,
	party *models.Party,
	categoryIDs []uuid.UUID,
	imageIDs []uuid.UUID,
	ticketCategories []models.TicketCategory,
) error {

	return s.tx.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			if err := s.parties.Update(
				ctx,
				party,
			); err != nil {
				return err
			}

			if err := s.partyCategories.Replace(
				tx,
				party.ID,
				categoryIDs,
			); err != nil {
				return err
			}

			if err := s.images.Replace(
				tx,
				party.ID,
				imageIDs,
			); err != nil {
				return err
			}

			if err := s.ticketCategories.Replace(
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

func (s *PartyCRUDService) CreateRelations(
	ctx context.Context,
	party *models.Party,
	categoryIDs []uuid.UUID,
	imageIDs []uuid.UUID,
) error {

	return s.tx.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			if err := s.parties.Create(
				tx,
				party,
			); err != nil {
				return err
			}

			// save categories in party_categories
			if err := s.partyCategories.Replace(
				tx,
				party.ID,
				categoryIDs,
			); err != nil {
				return err
			}

			// save images in party_media
			if err := s.images.Replace(
				tx,
				party.ID,
				imageIDs,
			); err != nil {
				return err
			}

			// create organizer membership
			member := &models.PartyMember{
				ID: uuid.New(),

				UserID: party.OrganizerID,

				PartyID: party.ID,
			}

			if err := s.members.Create(
				tx,
				member,
			); err != nil {
				return err
			}

			// assign organizer role
			role := &models.PartyMemberRole{
				ID: uuid.New(),

				PartyMemberID: member.ID,

				Role: enum.RoleOrganizer,
			}

			return s.roles.Create(
				tx,
				role,
			)
		},
	)
}
