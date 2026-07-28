package security

import (
	"errors"
	"regexp"
	"strings"
)

var usernameRegex = regexp.MustCompile(
	`^[a-zA-Z0-9_]+$`,
)

var ErrInvalidUsername = errors.New(
	"invalid username",
)

func ValidateUsername(
	username string,
) error {

	username = strings.TrimSpace(username)

	// length validation
	if len(username) < 3 || len(username) > 30 {
		return ErrInvalidUsername
	}

	// allowed characters
	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}

	return nil
}
