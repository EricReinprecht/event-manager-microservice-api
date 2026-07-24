package appErrors

import "errors"

var (
	ErrPartyMemberAlreadyExists   = errors.New("user is already a member of this party")
	ErrPartyNotFound              = errors.New("party not found")
	ErrTicketCategoryNotFound     = errors.New("ticket category not found")
	ErrNotEnoughTickets           = errors.New("not enough tickets available")
	ErrTicketAlreadyUsed          = errors.New("ticket already used")
	ErrNotAllowed                 = errors.New("not allowed")
	ErrInvalidPartyMemberRole     = errors.New("invalid party member role")
	ErrCannotRemoveOrganizer      = errors.New("party organizer cannot be removed")
	ErrPartyMemberNotFound        = errors.New("party member not found")
	ErrTicketNotValidNow          = errors.New("ticket is not valid at this time")
	ErrTicketAlreadyScanned       = errors.New("ticket was already scanned")
	ErrTicketAccessWindowRequired = errors.New("ticket category requires at least one access window")
	ErrTicketScanAlreadyDecided   = errors.New("ticket scan already decided")
)
