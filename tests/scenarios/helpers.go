package scenarios

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"gorm.io/gorm"
)

func CreatePurchase(
	t *testing.T,
	db *gorm.DB,
	userID uuid.UUID,
	partyID uuid.UUID,
) models.Purchase {

	t.Helper()

	purchase := models.Purchase{
		ID: uuid.New(),

		UserID: userID,

		PartyID: partyID,

		Status: enum.PurchaseStatusPaid,

		ExpiresAt: time.Now().
			UTC().
			Add(time.Hour),

		TotalPrice: 0,
	}

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	return purchase
}
