package http

import (
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
)

func ExecuteAuthenticatedRefund(
	router http.Handler,
	purchaseID uuid.UUID,
	userID uuid.UUID,
) *httptest.ResponseRecorder {

	req := AuthenticatedRequest(
		http.MethodPost,
		"/api/purchases/"+purchaseID.String()+"/refund",
		userID,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	return recorder
}
