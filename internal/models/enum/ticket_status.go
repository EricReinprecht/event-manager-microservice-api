package enum

type TicketStatus string

const (
	TicketStatusActive    TicketStatus = "ACTIVE"
	TicketStatusCancelled TicketStatus = "CANCELLED"
)
