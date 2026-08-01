package dto

import "github.com/google/uuid"

type CreateTicketCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`

	Price int64 `json:"price" binding:"min=0"`

	Capacity int `json:"capacity" binding:"required,min=1"`

	RequiresVerification bool `json:"requiresVerification"`

	RefundRequiresApproval bool `json:"refundRequiresApproval"`

	RefundPolicyID *uuid.UUID `json:"refundPolicyId"`

	AccessWindows []CreateAccessWindowRequest `json:"accessWindows" binding:"dive"`
}

type UpdateTicketCategoryRequest struct {
	ID *uuid.UUID `json:"id"`

	Name string `json:"name" binding:"required,min=2,max=100"`

	Price int64 `json:"price" binding:"min=0"`

	Capacity int `json:"capacity" binding:"required,min=1"`

	RequiresVerification bool `json:"requiresVerification"`

	RefundRequiresApproval bool `json:"refundRequiresApproval"`

	RefundPolicyID *uuid.UUID `json:"refundPolicyId"`

	AccessWindows []UpdateAccessWindowRequest `json:"accessWindows" binding:"dive"`
}

type TicketCategoryResponse struct {
	ID uuid.UUID `json:"id"`

	Name string `json:"name"`

	Price int64 `json:"price"`

	Capacity int `json:"capacity"`

	RequiresVerification bool `json:"requiresVerification"`

	RefundRequiresApproval bool `json:"refundRequiresApproval"`

	RefundPolicyID *uuid.UUID `json:"refundPolicyId"`

	AccessWindows []AccessWindowResponse `json:"accessWindows"`
}
