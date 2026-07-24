package fixtures

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
)

func User() models.User {

	return models.User{
		ID: uuid.New(),

		Email: "test-" + uuid.NewString() + "@example.com",

		Username: "user-" + uuid.NewString(),

		PasswordHash: "test-password-hash",

		FirstName: "Test",

		LastName: "User",
	}
}

func UserWithID(id uuid.UUID) models.User {

	return models.User{
		ID: id,

		Email: "test-" + uuid.NewString() + "@example.com",

		Username: "user-" + uuid.NewString(),

		PasswordHash: "test-password-hash",

		FirstName: "Test",

		LastName: "User",
	}
}
