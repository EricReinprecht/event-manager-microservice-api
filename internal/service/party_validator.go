package service

import (
	"context"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"

	"github.com/reinp/event-platform/backend/internal/repository"
	partyCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/party_category"
)

type PartyValidator struct {
	categories *partyCategoryRepository.Facade
	media      *repository.MediaRepository
}

func NewPartyValidator(
	categories *partyCategoryRepository.Facade,
	media *repository.MediaRepository,
) *PartyValidator {

	return &PartyValidator{
		categories: categories,
		media:      media,
	}
}

func (v *PartyValidator) Validate(
	ctx context.Context,
	categoryID uuid.UUID,
	thumbnailID *uuid.UUID,
	imageIDs []uuid.UUID,
) error {

	_, err := v.categories.Repository.FindByID(
		ctx,
		categoryID,
	)

	if err != nil {
		return appErrors.ErrCategoryNotFound
	}

	if thumbnailID != nil {

		_, err = v.media.FindByID(
			ctx,
			*thumbnailID,
		)

		if err != nil {
			return appErrors.ErrMediaNotFound
		}
	}

	for _, imageID := range imageIDs {

		_, err = v.media.FindByID(
			ctx,
			imageID,
		)

		if err != nil {
			return appErrors.ErrMediaNotFound
		}
	}

	return nil
}
