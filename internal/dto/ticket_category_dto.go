package dto

import "github.com/google/uuid"

type CreateTicketCategoryRequest struct {
	Name string `json:"name" binding:"required"`

	Price int64 `json:"price" binding:"required"`

	Capacity int `json:"capacity" binding:"required"`

	RequiresVerification bool `json:"requiresVerification"`

	RefundRequiresApproval bool `json:"refundRequiresApproval"`

	RefundPolicyID *uuid.UUID `json:"refundPolicyId"`

	AccessWindows []CreateAccessWindowRequest `json:"accessWindows"`
}

type UpdateTicketCategoryRequest struct {
	ID *uuid.UUID `json:"id"`

	Name string `json:"name"`

	Price int64 `json:"price"`

	Capacity int `json:"capacity"`

	RequiresVerification bool `json:"requiresVerification"`

	RefundRequiresApproval bool `json:"refundRequiresApproval"`

	RefundPolicyID *uuid.UUID `json:"refundPolicyId"`

	AccessWindows []UpdateAccessWindowRequest `json:"accessWindows"`
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
