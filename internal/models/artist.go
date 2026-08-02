package models

import "github.com/google/uuid"

type Artist struct {
	ID uuid.UUID

	UserID uuid.UUID

	StageName string

	Description string

	ImageID *uuid.UUID

	Verified bool
}
