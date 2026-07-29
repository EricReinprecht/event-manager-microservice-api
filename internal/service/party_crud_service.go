package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PartyCRUDService struct {
	parties    *repository.PartyRepository
	images     *repository.PartyImageRepository
	members    *repository.PartyMemberRepository
	roles      *repository.PartyMemberRoleRepository
	categories *repository.CategoryRepository
	media      *repository.MediaRepository

	tx *database.TransactionManager
}

func NewPartyCRUDService(
	parties *repository.PartyRepository,
	images *repository.PartyImageRepository,
	members *repository.PartyMemberRepository,
	roles *repository.PartyMemberRoleRepository,
	categories *repository.CategoryRepository,
	media *repository.MediaRepository,
	tx *database.TransactionManager,
) *PartyCRUDService {

	return &PartyCRUDService{
		parties:    parties,
		images:     images,
		members:    members,
		roles:      roles,
		categories: categories,
		media:      media,
		tx:         tx,
	}
}

func (s *PartyCRUDService) Create(
	ctx context.Context,
	party *models.Party,
	imageIDs []uuid.UUID,
) error {

	// Validate categories
	for _, category := range party.Categories {

		_, err := s.categories.FindByID(
			ctx,
			category.ID,
		)

		if err != nil {
			return appErrors.ErrCategoryNotFound
		}
	}

	// Validate thumbnail
	if party.ThumbnailID != nil {

		_, err := s.media.FindByID(
			ctx,
			*party.ThumbnailID,
		)

		if err != nil {
			return appErrors.ErrMediaNotFound
		}
	}

	// Validate images
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

			// Save many-to-many category relations
			if err := tx.
				Model(party).
				Association("Categories").
				Replace(party.Categories); err != nil {

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
				role)
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

func (s *PartyCRUDService) Delete(
	ctx context.Context,
	party *models.Party,
) error {

	return s.parties.Delete(
		ctx,
		party,
	)
}

func (s *PartyCRUDService) FindOwnedParty(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) (*models.Party, error) {

	party, err := s.FindByID(
		ctx,
		partyID,
	)

	if err != nil {
		return nil, err
	}

	if party.OrganizerID != userID {
		return nil, errors.New("not allowed")
	}

	return party, nil
}
