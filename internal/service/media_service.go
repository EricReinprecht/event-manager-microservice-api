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
	filename string,
	path string,
	mimeType string,
	size int64,
) (*models.Media, error) {
	media := &models.Media{
		Filename: filename,
		Path:     path,
		URL:      "/" + path,
		MimeType: mimeType,
		Size:     size,
	}

	if err := s.repository.Create(ctx, media); err != nil {
		return nil, err
	}
	return media, nil
}
