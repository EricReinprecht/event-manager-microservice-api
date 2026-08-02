package responses

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/i18n"
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

		translator := i18n.FromContext(c)

		translatedErrors := make(
			map[string]string,
			len(validationError.Errors),
		)

		for field, translationKey := range validationError.Errors {

			translatedErrors[field] =
				translator.T(
					translationKey,
				)
		}

		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"errors": translatedErrors,
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
	), errors.Is(
		err,
		appErrors.ErrForbidden,
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
