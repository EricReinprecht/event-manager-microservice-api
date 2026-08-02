package dto

import "github.com/google/uuid"

type CreatePurchaseRequest struct {
	Items []PurchaseItemRequest `json:"items" binding:"required,min=1"`
}

type PurchaseItemRequest struct {
	TicketCategoryID uuid.UUID `json:"ticket_category_id" binding:"required"`
	Quantity         int       `json:"quantity" binding:"required,min=1"`
}
