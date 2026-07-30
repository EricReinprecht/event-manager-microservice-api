package helpers

import (
	"errors"
	"time"
)

func ValidateParty(
	startAt time.Time,
	endAt time.Time,
	latitude float64,
	longitude float64,
	timezone string,
) error {

	if endAt.Before(startAt) {
		return errors.New("end date must be after start date")
	}

	if latitude == 0 || longitude == 0 {
		return errors.New("location coordinates are required")
	}

	if timezone == "" {
		return errors.New("timezone is required")
	}

	return nil
}
