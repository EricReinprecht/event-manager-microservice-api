package mapper

import (
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
)

func AccessWindowResponses(
	windows []models.TicketAccessWindow,
) []dto.AccessWindowResponse {

	result := make([]dto.AccessWindowResponse, 0, len(windows))

	for _, window := range windows {
		result = append(result, dto.AccessWindowResponse{
			ID:       window.ID,
			StartsAt: window.StartsAt,
			EndsAt:   window.EndsAt,
		})
	}

	return result
}
