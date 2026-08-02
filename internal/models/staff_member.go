package models

import "github.com/google/uuid"

type StaffMember struct {
	ID uuid.UUID

	PartyID uuid.UUID

	UserID *uuid.UUID

	FirstName string
	LastName  string

	Role        string
	Description string
}
