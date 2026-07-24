package enum

type TicketScanStatus string

const (
	TicketScanPending TicketScanStatus = "PENDING"

	TicketScanVerified TicketScanStatus = "VERIFIED"

	TicketScanRejected TicketScanStatus = "REJECTED"
)
