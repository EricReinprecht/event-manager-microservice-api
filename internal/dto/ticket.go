package dto

import "github.com/google/uuid"

type PurchaseTicketItem struct {
	TicketCategoryID uuid.UUID

	Quantity int
}
