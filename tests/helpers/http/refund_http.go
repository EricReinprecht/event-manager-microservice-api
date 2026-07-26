package http

import (
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
)

func ExecuteRefundRequest(
	router http.Handler,
	purchaseID uuid.UUID,
	token string,
) *httptest.ResponseRecorder {

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/purchases/"+purchaseID.String()+"/refund",
		nil,
	)

	if token != "" {

		req.Header.Set(
			"Authorization",
			"Bearer "+token,
		)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	return recorder
}
