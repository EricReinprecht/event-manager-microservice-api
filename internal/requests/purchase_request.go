package requests

import "github.com/google/uuid"

type PurchaseItemRequest struct {
	TicketCategoryID uuid.UUID `json:"ticket_category_id"`
	Quantity         int       `json:"quantity"`
}

type CreatePurchaseRequest struct {
	Items []PurchaseItemRequest `json:"items"`
}
