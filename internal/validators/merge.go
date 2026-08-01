package validators

import (
	"github.com/reinp/event-platform/backend/internal/appErrors"
)

func MergeValidationErrors(
	target appErrors.ValidationErrors,
	source appErrors.ValidationErrors,
) {

	for path, translationKey := range source {

		target[path] = translationKey
	}
}
