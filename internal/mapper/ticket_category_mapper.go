package mapper

import (
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
