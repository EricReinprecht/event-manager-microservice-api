package helpers

import "time"

func init() {
	time.Local = time.UTC
}
