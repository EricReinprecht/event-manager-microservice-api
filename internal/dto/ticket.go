package dto

import "github.com/google/uuid"

type PurchaseTicketItem struct {
	TicketCategoryID uuid.UUID `json:"ticket_category_id"`
	Quantity         int       `json:"quantity"`
}
