package fixtures

import (
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
)

func RefreshToken(
	userID uuid.UUID,
) models.RefreshToken {

	return models.RefreshToken{

		ID: uuid.New(),

		UserID: userID,

		FamilyID: uuid.New(),

		TokenHash: uuid.NewString(),

		UserAgent: "test-agent",

		IPAddress: "127.0.0.1",

		DeviceName: "test-device",

		ExpiresAt: time.Now().Add(
			30 * 24 * time.Hour,
		),
	}
}
