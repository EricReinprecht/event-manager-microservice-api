package helpers

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models"
)

func GetVerificationToken(
	db *gorm.DB,
	userID uuid.UUID,
) string {

	var verification models.EmailVerification

	err := db.
		Where(
			"user_id = ?",
			userID,
		).
		First(
			&verification,
		).
		Error

	if err != nil {
		panic(err)
	}

	// IMPORTANT:
	// database stores hashed token, not raw token.
	// If your tests need the raw token, store it when sending mail.
	return verification.Token
}
