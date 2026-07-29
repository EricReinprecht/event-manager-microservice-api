package service

import (
	"context"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PartyValidator struct {
	categoryRepository *repository.CategoryRepository
	mediaRepository    *repository.MediaRepository
}

func NewPartyValidator(
	categoryRepository *repository.CategoryRepository,
	mediaRepository *repository.MediaRepository,
) *PartyValidator {

	return &PartyValidator{
		categoryRepository: categoryRepository,
		mediaRepository:    mediaRepository,
	}
}

func (v *PartyValidator) Validate(
	ctx context.Context,
	categoryID uuid.UUID,
	thumbnailID *uuid.UUID,
	imageIDs []uuid.UUID,
) error {

	_, err := v.categoryRepository.FindByID(
		ctx,
		categoryID,
	)

	if err != nil {
		return appErrors.ErrCategoryNotFound
	}

	if thumbnailID != nil {

		_, err = v.mediaRepository.FindByID(
			ctx,
			*thumbnailID,
		)

		if err != nil {
			return appErrors.ErrMediaNotFound
		}
	}

	for _, imageID := range imageIDs {

		_, err = v.mediaRepository.FindByID(
			ctx,
			imageID,
		)

		if err != nil {
			return appErrors.ErrMediaNotFound
		}
	}

	return nil
}
