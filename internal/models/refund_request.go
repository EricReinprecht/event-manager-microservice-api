package models

import (
	"time"

	"github.com/google/uuid"
)

type RefundRequest struct {
	ID uuid.UUID

	PurchaseID uuid.UUID
	Purchase   Purchase

	RequestedBy uuid.UUID
	User        User

	Status enum.RefundRequestStatus

	ApprovedBy *uuid.UUID

	Reason string

	CreatedAt time.Time

	ApprovedAt *time.Time
}
