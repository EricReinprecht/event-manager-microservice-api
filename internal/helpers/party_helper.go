package helpers

import (
	"time"

	"github.com/reinp/event-platform/backend/internal/appErrors"
)

func ValidateParty(
	startAt time.Time,
	endAt time.Time,
	latitude float64,
	longitude float64,
	timezone string,
) error {

	validationErrors := appErrors.ValidationErrors{}

	if !endAt.After(startAt) {
		validationErrors["_section_schedule"] =
			appErrors.ErrMsgPartyEndBeforeStart
	}

	if latitude == 0 || longitude == 0 {
		validationErrors["location"] =
			appErrors.ErrMsgLocationRequired
	}

	if timezone == "" {
		validationErrors["timezone"] =
			appErrors.ErrMsgTimezoneRequired
	}

	if len(validationErrors) > 0 {
		return appErrors.NewValidationError(
			validationErrors,
		)
	}

	return nil
}
