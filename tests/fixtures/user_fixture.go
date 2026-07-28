package fixtures

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
)

func User() models.User {

	return models.User{
		ID: uuid.New(),

		Email: "test-" + uuid.NewString() + "@example.com",

		Username: "user-" + uuid.NewString(),

		PasswordHash: PasswordHash(
			"C4ctus!River#829Lamp",
		),

		FirstName: "Test",

		LastName: "User",
	}
}

func UserWithID(
	id uuid.UUID,
) models.User {

	return models.User{
		ID: id,

		Email: "test-" + uuid.NewString() + "@example.com",

		Username: "user-" + uuid.NewString(),

		PasswordHash: PasswordHash(
			"C4ctus!River#829Lamp",
		),

		FirstName: "Test",

		LastName: "User",
	}
}

func VerifiedUser() models.User {

	now := time.Now()

	user := User()

	user.VerifiedAt = &now

	return user
}

func VerifiedUserWithID(
	id uuid.UUID,
) models.User {

	now := time.Now()

	user := UserWithID(
		id,
	)

	user.VerifiedAt = &now

	return user
}

func PasswordHash(
	password string,
) string {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		panic(err)
	}

	return string(hash)
}
