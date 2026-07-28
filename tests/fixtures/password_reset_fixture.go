package fixtures

import (
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
)

func PasswordResetToken(
	userID uuid.UUID,
	hash string,
) models.PasswordResetToken {

	return models.PasswordResetToken{

		ID: uuid.New(),

		UserID: userID,

		TokenHash: hash,

		ExpiresAt: time.Now().Add(
			15 * time.Minute,
		),
	}
}
