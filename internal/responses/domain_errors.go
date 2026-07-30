package responses

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/appErrors"
)

func HandleDomainError(
	c *gin.Context,
	err error,
) {

	switch {

	case errors.Is(
		err,
		appErrors.ErrCategoryNotFound,
	),
		errors.Is(
			err,
			appErrors.ErrMediaNotFound,
		):

		Error(
			c,
			http.StatusBadRequest,
			err,
		)

	case errors.Is(
		err,
		appErrors.ErrUnauthorized,
	):

		Error(
			c,
			http.StatusForbidden,
			err,
		)

	case errors.Is(
		err,
		appErrors.ErrPartyNotFound,
	):

		Error(
			c,
			http.StatusNotFound,
			err,
		)

	default:

		InternalError(
			c,
			err,
		)
	}
}
