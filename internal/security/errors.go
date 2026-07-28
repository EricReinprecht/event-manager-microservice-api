package security

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func ErrorMessage(err error) string {

	var validationErr validator.ValidationErrors

	if errors.As(err, &validationErr) {

		field := validationErr[0]

		switch field.Tag() {

		case "required":
			return field.Field() + " is required"

		case "email":
			return field.Field() + " must be a valid email"

		case "min":
			return field.Field() + " is too short"

		case "max":
			return field.Field() + " is too long"

		default:
			return field.Field() + " is invalid"
		}
	}

	return err.Error()
}
