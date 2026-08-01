package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reinp/event-platform/backend/internal/clock"
)

type JWT struct {
	Secret              string
	Clock               clock.Clock
	accessTokenDuration time.Duration
}

func NewJWT(
	secret string,
	clock clock.Clock,
	accessTokenDuration time.Duration,

) *JWT {

	return &JWT{
		Secret:              secret,
		Clock:               clock,
		accessTokenDuration: accessTokenDuration,
	}
}

func (j *JWT) Generate(
	userID string,
	familyID string,
) (string, error) {

	claims := jwt.MapClaims{
		"user_id":   userID,
		"family_id": familyID,
		"exp": j.Clock.Now().
			Add(j.accessTokenDuration).
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
