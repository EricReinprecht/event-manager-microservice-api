package helpers

import (
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/models"
)

func CreateEmailVerification(
	userID uuid.UUID,
) (string, *models.EmailVerification) {

	rawToken := auth.GenerateToken()

	return rawToken, &models.EmailVerification{

		UserID: userID,

		Token: auth.HashToken(
			rawToken,
		),

		ExpiresAt: time.Now().Add(
			24 * time.Hour,
		),
	}
}
