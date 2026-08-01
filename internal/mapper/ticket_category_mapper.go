package mapper

import (
	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
)

func TicketCategoryResponse(
	category *models.TicketCategory,
) dto.TicketCategoryResponse {

	return dto.TicketCategoryResponse{
		ID: category.ID,

		Name: category.Name,

		Price: category.Price,

		Capacity: category.Capacity,

		RequiresVerification: category.RequiresVerification,

		RefundRequiresApproval: category.RefundRequiresApproval,

		RefundPolicyID: category.RefundPolicyID,

		AccessWindows: AccessWindowResponses(
			category.AccessWindows,
		),
	}
}

func TicketCategoryResponses(
	categories []models.TicketCategory,
) []dto.TicketCategoryResponse {

	result := make(
		[]dto.TicketCategoryResponse,
		0,
		len(categories),
	)

	for _, category := range categories {

		result = append(
			result,
			dto.TicketCategoryResponse{
				ID: category.ID,

				Name: category.Name,

				Price: category.Price,

				Capacity: category.Capacity,

				RequiresVerification: category.RequiresVerification,

				RefundRequiresApproval: category.RefundRequiresApproval,

				RefundPolicyID: category.RefundPolicyID,

				AccessWindows: AccessWindowResponses(
					category.AccessWindows,
				),
			},
		)
	}

	return result
}

func TicketCategory(
	category dto.UpdateTicketCategoryRequest,
	partyID uuid.UUID,
) models.TicketCategory {

	result := models.TicketCategory{
		PartyID:                partyID,
		Name:                   category.Name,
		Price:                  category.Price,
		Capacity:               category.Capacity,
		RequiresVerification:   category.RequiresVerification,
		RefundRequiresApproval: category.RefundRequiresApproval,
		RefundPolicyID:         category.RefundPolicyID,
	}

	if category.ID != nil {
		result.ID = *category.ID
	}

	result.AccessWindows = AccessWindows(
		category.AccessWindows,
		result.ID,
	)

	return result
}

func TicketCategories(
	categories []dto.UpdateTicketCategoryRequest,
	partyID uuid.UUID,
) []models.TicketCategory {

	result := make(
		[]models.TicketCategory,
		0,
		len(categories),
	)

	for _, category := range categories {

		result = append(
			result,
			TicketCategory(
				category,
				partyID,
			),
		)
	}

	return result
}

func AccessWindow(
	window dto.UpdateAccessWindowRequest,
	ticketCategoryID uuid.UUID,
) models.TicketAccessWindow {

	result := models.TicketAccessWindow{
		TicketCategoryID: ticketCategoryID,
		StartsAt:         window.StartsAt,
		EndsAt:           window.EndsAt,
	}

	if window.ID != nil {
		result.ID = *window.ID
	}

	return result
}

func AccessWindows(
	windows []dto.UpdateAccessWindowRequest,
	ticketCategoryID uuid.UUID,
) []models.TicketAccessWindow {

	result := make(
		[]models.TicketAccessWindow,
		0,
		len(windows),
	)

	for _, window := range windows {

		result = append(
			result,
			AccessWindow(
				window,
				ticketCategoryID,
			),
		)
	}

	return result
}

func AccessWindowResponses(
	windows []models.TicketAccessWindow,
) []dto.AccessWindowResponse {

	result := make(
		[]dto.AccessWindowResponse,
		0,
		len(windows),
	)

	for _, window := range windows {

		result = append(
			result,
			dto.AccessWindowResponse{
				ID:       window.ID,
				StartsAt: window.StartsAt,
				EndsAt:   window.EndsAt,
			},
		)
	}

	return result
}
