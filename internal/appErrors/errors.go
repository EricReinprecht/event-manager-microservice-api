package appErrors

import "errors"

var (
	ErrPartyMemberAlreadyExists = errors.New("user is already a member of this party")
)
