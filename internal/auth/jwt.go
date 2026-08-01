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

func (j *JWT) Generate(
	userID string,
	familyID string,
) (string, error) {

	claims := jwt.MapClaims{
		"user_id":   userID,
		"family_id": familyID,
		"exp": time.Now().Add(
			15 * time.Second,
		).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(j.Secret),
	)
}
