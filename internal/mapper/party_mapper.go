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
		Title:        req.Title,
		Description:  req.Description,
		LocationName: req.LocationName,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Timezone:     req.Timezone,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
		OrganizerID:  organizerID,
		ThumbnailID:  req.ThumbnailID,
		Categories:   []models.Category{},
	}
}

func ApplyPartyUpdate(
	party *models.Party,
	req dto.UpdatePartyRequest,
) {
	party.Title = req.Title
	party.Description = req.Description
	party.LocationName = req.LocationName
	party.Latitude = req.Latitude
	party.Longitude = req.Longitude
	party.Timezone = req.Timezone
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

		Latitude: party.Latitude,

		Longitude: party.Longitude,

		Timezone: party.Timezone,

		StartAt: party.StartAt,

		EndAt: party.EndAt,

		ThumbnailID: party.ThumbnailID,

		OrganizerID: party.OrganizerID,

		Categories: CategoryResponses(party.Categories),
	}
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
	categories []models.Category,
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
