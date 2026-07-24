package helpers

import "time"

type Clock interface {
	Now() time.Time
}
