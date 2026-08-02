package models

import (
	"time"

	"github.com/google/uuid"
)

type ArtistSlot struct {
	ID uuid.UUID

	PartyID uuid.UUID

	ArtistID uuid.UUID

	StageID uuid.UUID

	StartAt time.Time

	EndAt time.Time
}
