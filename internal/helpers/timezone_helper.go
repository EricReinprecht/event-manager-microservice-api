package helpers

import "time"

func DayRangeUTC(
	date string,
	timezone string,
) (time.Time, time.Time, error) {

	loc, err := time.LoadLocation(timezone)

	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	startLocal, err := time.ParseInLocation(
		"2006-01-02",
		date,
		loc,
	)

	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endLocal := startLocal.AddDate(0, 0, 1)

	return startLocal.UTC(), endLocal.UTC(), nil
}
