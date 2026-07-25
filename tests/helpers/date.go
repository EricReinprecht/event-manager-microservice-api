package helpers

import "time"

func UTCDate(
	year int,
	month time.Month,
	day int,
	hour int,
) time.Time {

	return time.Date(
		year,
		month,
		day,
		hour,
		0,
		0,
		0,
		time.UTC,
	)
}
