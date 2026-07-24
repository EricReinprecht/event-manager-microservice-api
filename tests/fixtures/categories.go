package fixtures

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
)

func Category() models.Category {

	return models.Category{

		ID: uuid.New(),

		Name: "Festival",
	}
}

func CategoryWithID(
	id uuid.UUID,
) models.Category {

	return models.Category{

		ID: id,

		Name: "Festival",
	}
}
