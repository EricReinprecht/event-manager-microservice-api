package helpers

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/reinp/event-platform/backend/internal/appErrors"
)

func BindingValidationErrors(
	err error,
	request any,
) *appErrors.ValidationError {

	var validationErrors validator.ValidationErrors

	if !errors.As(
		err,
		&validationErrors,
	) {
		return nil
	}

	requestType := reflect.TypeOf(
		request,
	)

	if requestType == nil {
		return nil
	}

	if requestType.Kind() == reflect.Pointer {
		requestType = requestType.Elem()
	}

	result := appErrors.ValidationErrors{}

	for _, fieldError := range validationErrors {

		path := validationFieldPath(
			fieldError,
			requestType,
		)

		if path == "" {
			continue
		}

		result[path] =
			validationMessageKey(
				fieldError.Tag(),
			)
	}

	if len(result) == 0 {
		return nil
	}

	return appErrors.NewValidationError(
		result,
	)
}

func validationFieldPath(
	fieldError validator.FieldError,
	rootType reflect.Type,
) string {

	namespace :=
		fieldError.StructNamespace()

	parts := strings.Split(
		namespace,
		".",
	)

	if len(parts) > 0 {
		parts = parts[1:]
	}

	currentType := rootType

	pathParts := make(
		[]string,
		0,
		len(parts),
	)

	for _, part := range parts {

		fieldName, indexes :=
			splitValidationPart(part)

		for currentType.Kind() ==
			reflect.Pointer {

			currentType =
				currentType.Elem()
		}

		if currentType.Kind() !=
			reflect.Struct {

			pathParts = append(
				pathParts,
				lowerFirst(fieldName)+indexes,
			)

			continue
		}

		structField, found :=
			currentType.FieldByName(
				fieldName,
			)

		if !found {

			pathParts = append(
				pathParts,
				lowerFirst(fieldName)+indexes,
			)

			continue
		}

		pathParts = append(
			pathParts,
			jsonFieldName(structField)+indexes,
		)

		currentType =
			structField.Type

		for currentType.Kind() ==
			reflect.Pointer {

			currentType =
				currentType.Elem()
		}

		if currentType.Kind() ==
			reflect.Slice ||
			currentType.Kind() ==
				reflect.Array {

			currentType =
				currentType.Elem()
		}
	}

	return strings.Join(
		pathParts,
		".",
	)
}

func splitValidationPart(
	part string,
) (string, string) {

	index := strings.Index(
		part,
		"[",
	)

	if index == -1 {
		return part, ""
	}

	fieldName := part[:index]

	indexes := part[index:]

	indexes = strings.ReplaceAll(
		indexes,
		"][",
		".",
	)

	indexes = strings.ReplaceAll(
		indexes,
		"[",
		".",
	)

	indexes = strings.ReplaceAll(
		indexes,
		"]",
		"",
	)

	return fieldName, indexes
}

func jsonFieldName(
	field reflect.StructField,
) string {

	jsonTag := field.Tag.Get(
		"json",
	)

	if jsonTag == "" {
		return lowerFirst(
			field.Name,
		)
	}

	name := strings.Split(
		jsonTag,
		",",
	)[0]

	if name == "" || name == "-" {
		return lowerFirst(
			field.Name,
		)
	}

	return name
}

func lowerFirst(
	value string,
) string {

	if value == "" {
		return value
	}

	return strings.ToLower(
		value[:1],
	) + value[1:]
}

func validationMessageKey(
	tag string,
) string {

	switch tag {

	case "required":
		return "validation.required"

	case "min":
		return "validation.min"

	case "max":
		return "validation.max"

	default:
		return "validation.invalid"
	}
}
