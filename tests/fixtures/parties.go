package fixtures

import (
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
)

func Party() models.Party {

	return models.Party{
		ID: uuid.New(),

		Title: "Test Festival",

		Description: "Test event",

		StartAt: time.Now().
			Add(24 * time.Hour),

		EndAt: time.Now().
			Add(48 * time.Hour),

		Location: "Test Location",

		CategoryID: uuid.Nil,

		OrganizerID: uuid.New(),
	}
}

func PartyWithOrganizer(
	organizerID uuid.UUID,
) models.Party {

	return models.Party{

		ID: uuid.New(),

		Title: "Test Festival",

		Description: "Test event",

		StartAt: time.Now().
			Add(24 * time.Hour),

		EndAt: time.Now().
			Add(48 * time.Hour),

		Location: "Test Location",

		OrganizerID: organizerID,
	}
}
