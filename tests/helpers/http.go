package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func JSONRequest(
	method string,
	url string,
	body any,
) *http.Request {

	var data []byte

	if body != nil {

		jsonBody, err := json.Marshal(body)

		if err != nil {
			panic(err)
		}

		data = jsonBody
	}

	req := httptest.NewRequest(
		method,
		url,
		bytes.NewBuffer(data),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	return req
}

func AuthenticatedRequest(
	method string,
	url string,
	body any,
	token string,
) *http.Request {

	req := JSONRequest(
		method,
		url,
		body,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	return req
}

func ExecuteRequest(
	router *gin.Engine,
	req *http.Request,
) *httptest.ResponseRecorder {

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	return recorder
}

func DecodeJSON(
	recorder *httptest.ResponseRecorder,
	target any,
) {

	err := json.Unmarshal(
		recorder.Body.Bytes(),
		target,
	)

	if err != nil {
		panic(err)
	}
}
