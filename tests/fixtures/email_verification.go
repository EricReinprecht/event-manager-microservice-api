package fixtures

import (
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
)

func EmailVerification(
	userID uuid.UUID,
	hash string,
) models.EmailVerification {

	return models.EmailVerification{

		ID: uuid.New(),

		UserID: userID,

		Token: hash,

		ExpiresAt: time.Now().Add(
			24 * time.Hour,
		),
	}
}
