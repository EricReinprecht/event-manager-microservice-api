package helpers

import (
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/auth"
)

func CreateAuthToken(
	userID uuid.UUID,
) string {

	jwtService := auth.NewJWT(
		"test-secret",
		NewFakeClock(
			time.Now().UTC(),
		),
		15*time.Minute,
	)

	token, err := jwtService.Generate(
		userID.String(),
		uuid.New().String(),
	)

	if err != nil {
		panic(err)
	}

	return token
}
