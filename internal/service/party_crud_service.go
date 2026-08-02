package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"

	partyRepository "github.com/reinp/event-platform/backend/internal/repository/party"
	partyCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/party_category"
)

type PartyCRUDService struct {
	parties    *partyRepository.Facade
	categories *partyCategoryRepository.Facade
	media      *repository.MediaRepository
}

func (s *PartyCRUDService) Publish(
	ctx context.Context,
	partyID uuid.UUID,
	publishedAt time.Time,
) error {
	return s.parties.Repository.Publish(ctx, partyID, publishedAt)
}

func NewPartyCRUDService(
	parties *partyRepository.Facade,
	categories *partyCategoryRepository.Facade,
	media *repository.MediaRepository,
) *PartyCRUDService {

	return &PartyCRUDService{
		parties:    parties,
		categories: categories,
		media:      media,
	}
}

func (s *PartyCRUDService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Party, error) {

	return s.parties.Repository.FindByID(
		ctx,
		id,
	)
}

func (s *PartyCRUDService) Delete(
	ctx context.Context,
	party *models.Party,
) error {

	return s.parties.Repository.Delete(
		ctx,
		party,
	)
}

func (s *PartyCRUDService) CreateRelations(
	ctx context.Context,
	party *models.Party,
	categoryIDs []uuid.UUID,
	imageIDs []uuid.UUID,
	ticketCategories []models.TicketCategory,
) error {

	if err := s.validateCategories(
		ctx,
		categoryIDs,
	); err != nil {

		return err
	}

	if err := s.validateMedia(
		ctx,
		party.ThumbnailID,
		imageIDs,
	); err != nil {

		return err
	}

	return s.parties.Write.CreateWithRelations(
		ctx,
		party,
		categoryIDs,
		imageIDs,
		ticketCategories,
	)
}

func (s *PartyCRUDService) UpdateRelations(
	ctx context.Context,
	party *models.Party,
	categoryIDs []uuid.UUID,
	imageIDs []uuid.UUID,
	ticketCategories []models.TicketCategory,
) error {

	if err := s.validateCategories(
		ctx,
		categoryIDs,
	); err != nil {

		return err
	}

	if err := s.validateMedia(
		ctx,
		party.ThumbnailID,
		imageIDs,
	); err != nil {

		return err
	}

	return s.parties.Write.UpdateWithRelations(
		ctx,
		party,
		categoryIDs,
		imageIDs,
		ticketCategories,
	)
}

func (s *PartyCRUDService) UpdateImages(
	ctx context.Context,
	partyID uuid.UUID,
	imageIDs []uuid.UUID,
) error {

	if err := s.validateMedia(
		ctx,
		nil,
		imageIDs,
	); err != nil {

		return err
	}

	return s.parties.Write.ReplaceImages(
		ctx,
		partyID,
		imageIDs,
	)
}

func (s *PartyCRUDService) UpdateCategories(
	ctx context.Context,
	partyID uuid.UUID,
	categoryIDs []uuid.UUID,
) error {

	if err := s.validateCategories(
		ctx,
		categoryIDs,
	); err != nil {

		return err
	}

	return s.parties.Write.ReplaceCategories(
		ctx,
		partyID,
		categoryIDs,
	)
}

func (s *PartyCRUDService) validateCategories(
	ctx context.Context,
	categoryIDs []uuid.UUID,
) error {

	for _, categoryID := range categoryIDs {

		if _, err := s.categories.Repository.FindByID(
			ctx,
			categoryID,
		); err != nil {

			return appErrors.ErrCategoryNotFound
		}
	}

	return nil
}

func (s *PartyCRUDService) validateMedia(
	ctx context.Context,
	thumbnailID *uuid.UUID,
	imageIDs []uuid.UUID,
) error {

	if thumbnailID != nil {

		if _, err := s.media.FindByID(
			ctx,
			*thumbnailID,
		); err != nil {

			return appErrors.ErrMediaNotFound
		}
	}

	for _, imageID := range imageIDs {

		if _, err := s.media.FindByID(
			ctx,
			imageID,
		); err != nil {

			return appErrors.ErrMediaNotFound
		}
	}

	return nil
}
