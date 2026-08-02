package mapper

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
)

func NewParty(
	req dto.CreatePartyRequest,
	organizerID uuid.UUID,
) *models.Party {

	return &models.Party{
		Title:       req.Title,
		Description: req.Description,

		LocationName: req.LocationName,

		Street:      req.Location.Street,
		HouseNumber: req.Location.HouseNumber,
		City:        req.Location.City,
		Country:     req.Location.Country,
		PostalCode:  req.Location.PostalCode,

		Latitude:  req.Location.Latitude,
		Longitude: req.Location.Longitude,
		Timezone:  req.Location.Timezone,

		StartAt: req.StartAt,
		EndAt:   req.EndAt,

		OrganizerID: organizerID,

		ThumbnailID: req.ThumbnailID,

		Categories: []models.PartyCategory{},
	}
}

func ApplyPartyUpdate(
	party *models.Party,
	req dto.UpdatePartyRequest,
) {
	party.Title = req.Title

	party.Description = req.Description

	party.LocationName = req.LocationName

	// Location metadata
	party.Street = req.Location.Street
	party.HouseNumber = req.Location.HouseNumber
	party.City = req.Location.City
	party.Country = req.Location.Country
	party.PostalCode = req.Location.PostalCode

	party.Latitude = req.Location.Latitude
	party.Longitude = req.Location.Longitude
	party.Timezone = req.Location.Timezone

	party.StartAt = req.StartAt

	party.EndAt = req.EndAt

	party.ThumbnailID = req.ThumbnailID
}

func PartyResponse(
	party *models.Party,
) dto.PartyResponse {

	return dto.PartyResponse{
		ID: party.ID,

		Title: party.Title,

		Description: party.Description,

		LocationName: party.LocationName,

		Location: dto.PartyLocation{
			Street: party.Street,

			HouseNumber: party.HouseNumber,

			City: party.City,

			Country: party.Country,

			PostalCode: party.PostalCode,

			Latitude: party.Latitude,

			Longitude: party.Longitude,

			Timezone: party.Timezone,
		},

		StartAt: party.StartAt,

		EndAt: party.EndAt,

		ThumbnailID: party.ThumbnailID,

		ImageIDs: mediaIDs(party.Images),

		Thumbnail: MediaResponse(party.Thumbnail),

		Images: MediaResponses(party.Images),

		OrganizerID: party.OrganizerID,

		Categories: CategoryResponses(party.Categories),

		TicketCategories: TicketCategoryResponses(
			party.TicketCategories,
		),
	}
}

func MediaResponse(media *models.Media) *dto.MediaResponse {
	if media == nil {
		return nil
	}
	return &dto.MediaResponse{ID: media.ID, Filename: media.Filename, URL: media.URL, MimeType: media.MimeType}
}

func MediaResponses(media []models.Media) []dto.MediaResponse {
	result := make([]dto.MediaResponse, 0, len(media))
	for index := range media {
		result = append(result, *MediaResponse(&media[index]))
	}
	return result
}

func mediaIDs(media []models.Media) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(media))
	for _, item := range media {
		ids = append(ids, item.ID)
	}
	return ids
}

func PartyResponses(
	parties []models.Party,
) []dto.PartyResponse {

	result := make([]dto.PartyResponse, 0, len(parties))

	for _, party := range parties {

		result = append(
			result,
			PartyResponse(&party),
		)
	}

	return result
}

func CategoryResponses(
	categories []models.PartyCategory,
) []dto.CategoryResponse {

	result := make([]dto.CategoryResponse, 0, len(categories))

	for _, category := range categories {
		result = append(result, dto.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	return result
}
