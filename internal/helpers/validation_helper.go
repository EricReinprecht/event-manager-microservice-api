package helpers

import "github.com/reinp/event-platform/backend/internal/appErrors"

func MergeValidationErrors(
	target appErrors.ValidationErrors,
	source appErrors.ValidationErrors,
) {

	for path, message := range source {
		target[path] = message
	}
}
