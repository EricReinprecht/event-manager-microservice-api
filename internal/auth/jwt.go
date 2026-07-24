package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reinp/event-platform/backend/internal/clock"
)

type JWT struct {
	Secret string
	Clock  clock.Clock
}

func NewJWT(
	secret string,
	clock clock.Clock,
) *JWT {

	return &JWT{
		Secret: secret,
		Clock:  clock,
	}
}

func (j *JWT) Generate(userID string) (string, error) {

	now := j.Clock.Now()

	claims := jwt.MapClaims{
		"user_id": userID,

		"exp": now.
			Add(24 * time.Hour).
			Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(j.Secret),
	)
}
