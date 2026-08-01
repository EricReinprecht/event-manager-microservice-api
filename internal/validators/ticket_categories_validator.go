package validators

import (
	"fmt"
	"strings"
	"time"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/dto"
)

type ticketCategoryValidationItem struct {
	Name          string
	Price         int64
	Capacity      int
	AccessWindows []accessWindowValidationItem
}

type accessWindowValidationItem struct {
	StartAt time.Time
	EndAt   time.Time
}

func ValidateCreateTicketCategories(
	categories []dto.CreateTicketCategoryRequest,
) appErrors.ValidationErrors {

	items := make(
		[]ticketCategoryValidationItem,
		0,
		len(categories),
	)

	for _, category := range categories {

		windows := make(
			[]accessWindowValidationItem,
			0,
			len(category.AccessWindows),
		)

		for _, window := range category.AccessWindows {

			windows = append(
				windows,
				accessWindowValidationItem{
					StartAt: window.StartsAt,
					EndAt:   window.EndsAt,
				},
			)
		}

		items = append(
			items,
			ticketCategoryValidationItem{
				Name:          category.Name,
				Price:         category.Price,
				Capacity:      category.Capacity,
				AccessWindows: windows,
			},
		)
	}

	return validateTicketCategories(items)
}

func ValidateUpdateTicketCategories(
	categories []dto.UpdateTicketCategoryRequest,
) appErrors.ValidationErrors {

	items := make(
		[]ticketCategoryValidationItem,
		0,
		len(categories),
	)

	for _, category := range categories {

		windows := make(
			[]accessWindowValidationItem,
			0,
			len(category.AccessWindows),
		)

		for _, window := range category.AccessWindows {

			windows = append(
				windows,
				accessWindowValidationItem{
					StartAt: window.StartsAt,
					EndAt:   window.EndsAt,
				},
			)
		}

		items = append(
			items,
			ticketCategoryValidationItem{
				Name:          category.Name,
				Price:         category.Price,
				Capacity:      category.Capacity,
				AccessWindows: windows,
			},
		)
	}

	return validateTicketCategories(items)
}

func validateTicketCategories(
	categories []ticketCategoryValidationItem,
) appErrors.ValidationErrors {

	validationErrors :=
		appErrors.ValidationErrors{}

	seenNames := make(
		map[string]int,
	)

	for categoryIndex, category := range categories {

		name := strings.TrimSpace(
			category.Name,
		)

		namePath := fmt.Sprintf(
			"ticketCategories.%d.name",
			categoryIndex,
		)

		if name == "" {

			validationErrors[namePath] =
				appErrors.ErrMsgTicketCategoryNameRequired

		} else {

			normalizedName :=
				strings.ToLower(name)

			if _, exists :=
				seenNames[normalizedName]; exists {

				validationErrors[namePath] =
					appErrors.ErrMsgTicketCategoryNameDuplicate

			} else {

				seenNames[normalizedName] =
					categoryIndex
			}
		}

		if category.Price < 0 {

			validationErrors[fmt.Sprintf(
				"ticketCategories.%d.price",
				categoryIndex,
			)] =
				appErrors.ErrMsgTicketCategoryPriceInvalid
		}

		if category.Capacity < 1 {

			validationErrors[fmt.Sprintf(
				"ticketCategories.%d.capacity",
				categoryIndex,
			)] =
				appErrors.ErrMsgTicketCategoryCapacityInvalid
		}

		for windowIndex, window := range category.AccessWindows {

			if !window.EndAt.After(window.StartAt) {

				validationErrors[fmt.Sprintf(
					"ticketCategories.%d.accessWindows.%d._repeater",
					categoryIndex,
					windowIndex,
				)] =
					appErrors.ErrMsgAccessWindowEndBeforeStart
			}
		}
	}

	return validationErrors
}
