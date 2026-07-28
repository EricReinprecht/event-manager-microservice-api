package helpers

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func DoAuthenticatedRequest(
	t *testing.T,
	router *gin.Engine,
	method string,
	url string,
	token string,
	body any,
) HTTPJSONResponse {

	req := AuthenticatedRequest(
		method,
		url,
		body,
		token,
	)

	recorder := ExecuteRequest(
		router,
		req,
	)

	var response any

	DecodeJSON(
		recorder,
		&response,
	)

	return HTTPJSONResponse{
		StatusCode: recorder.Code,
		Body:       response,
	}
}
