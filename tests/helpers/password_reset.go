package helpers

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

func CreatePasswordResetToken(
	db database.DBExecutor,
	userID uuid.UUID,
) (string, error) {

	ctx := context.Background()

	rawToken := auth.GenerateToken()

	reset := &models.PasswordResetToken{
		UserID: userID,

		TokenHash: auth.HashToken(
			rawToken,
		),

		ExpiresAt: time.Now().Add(
			15 * time.Minute,
		),
	}

	repo := repository.NewPasswordResetTokenRepository(
		db,
	)

	err := repo.Create(
		ctx,
		reset,
	)

	if err != nil {
		return "", err
	}

	return rawToken, nil
}
