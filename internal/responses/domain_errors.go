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

	var validationError *appErrors.ValidationError

	if errors.As(
		err,
		&validationError,
	) {

		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"errors": validationError.Errors,
			},
		)

		return
	}

	switch {

	case errors.Is(
		err,
		appErrors.ErrCategoryNotFound,
	),
		errors.Is(
			err,
			appErrors.ErrMediaNotFound,
		),
		errors.Is(
			err,
			appErrors.ErrTicketCategoryExists,
		),
		errors.Is(
			err,
			appErrors.ErrTicketAccessWindowRequired,
		),
		errors.Is(
			err,
			appErrors.ErrAccessWindowInvalid,
		):

		BadRequest(
			c,
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

		NotFound(
			c,
			err,
		)

	default:

		InternalError(
			c,
			err,
		)
	}
}
