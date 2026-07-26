package http

import (
	"net/http"
	"net/http/httptest"
)

func ExecuteRequest(
	router http.Handler,
	method string,
	url string,
	token string,
) *httptest.ResponseRecorder {

	req := httptest.NewRequest(
		method,
		url,
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
