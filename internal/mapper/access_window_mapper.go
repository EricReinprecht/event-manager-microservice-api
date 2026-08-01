package mapper

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
)

func AccessWindowsFromRequest(
	windows []dto.CreateAccessWindowRequest,
) []models.TicketAccessWindow {

	result := make(
		[]models.TicketAccessWindow,
		0,
		len(windows),
	)

	for _, window := range windows {

		result = append(
			result,
			models.TicketAccessWindow{
				StartsAt: window.StartsAt,
				EndsAt:   window.EndsAt,
			},
		)
	}

	return result
}

func AccessWindowsFromUpdateRequest(
	windows []dto.UpdateAccessWindowRequest,
) []models.TicketAccessWindow {

	result := make(
		[]models.TicketAccessWindow,
		0,
		len(windows),
	)

	for _, window := range windows {

		result = append(
			result,
			models.TicketAccessWindow{
				ID:       uuidOrNil(window.ID),
				StartsAt: window.StartsAt,
				EndsAt:   window.EndsAt,
			},
		)
	}

	return result
}

func uuidOrNil(id *uuid.UUID) uuid.UUID {

	if id == nil {
		return uuid.Nil
	}

	return *id
}
