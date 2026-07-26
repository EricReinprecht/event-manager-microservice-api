package helpers

import (
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/auth"
)

func CreateAuthToken(
	userID uuid.UUID,
) string {

	jwt := auth.NewJWT(
		"test-secret",
		NewFakeClock(
			time.Now().UTC(),
		),
	)

	token, err := jwt.Generate(
		userID.String(),
	)

	if err != nil {
		panic(err)
	}

	return token
}
