package models

import "github.com/google/uuid"

type PartyStage struct {
	ID uuid.UUID

	PartyID uuid.UUID

	Name string

	Description string
}
