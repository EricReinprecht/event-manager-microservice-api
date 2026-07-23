package appErrors

import "errors"

var (
	ErrPartyMemberAlreadyExists = errors.New("user is already a member of this party")
	ErrPartyNotFound            = errors.New("party not found")
	ErrTicketCategoryNotFound   = errors.New("ticket category not found")
	ErrNotEnoughTickets         = errors.New("not enough tickets available")
)
