package http

import (
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func AuthenticatedRequest(
	method string,
	url string,
	userID uuid.UUID,
) *http.Request {

	req := httptest.NewRequest(
		method,
		url,
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+helpers.CreateAuthToken(userID),
	)

	return req
}
