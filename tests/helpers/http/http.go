package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func AuthHeader(
	userID uuid.UUID,
) http.Header {

	token := helpers.CreateAuthToken(
		userID,
	)

	headers := http.Header{}

	headers.Set(
		"Authorization",
		"Bearer "+token,
	)

	return headers
}
