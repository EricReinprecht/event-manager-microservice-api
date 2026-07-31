package mapper

import (
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
)

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
