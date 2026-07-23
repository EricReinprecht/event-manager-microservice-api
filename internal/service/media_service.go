package service

import (
	"context"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type MediaService struct {
	repository *repository.MediaRepository
}

func NewMediaService(
	repository *repository.MediaRepository,
) *MediaService {

	return &MediaService{
		repository: repository,
	}
}

func (s *MediaService) Create(
	ctx context.Context,
	media *models.Media,
) error {

	return s.repository.Create(
		ctx,
		media,
	)
}
