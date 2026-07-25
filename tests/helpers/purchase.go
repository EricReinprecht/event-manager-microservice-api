package helpers

import (
	"testing"

	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"
	"github.com/reinp/event-platform/backend/tests/fixtures"
)

func NewPurchaseService(
	db *gorm.DB,
) *service.PurchaseService {

	executor := database.NewGormExecutor(db)

	repository := repository.NewPurchaseRepository(
		executor,
	)

	return service.NewPurchaseService(
		repository,
	)
}

func CreatePurchase(
	t *testing.T,
	db *gorm.DB,
	user *models.User,
	party *models.Party,
	status enum.PurchaseStatus,
) *models.Purchase {

	t.Helper()

	purchase := fixtures.Purchase(
		user.ID,
		party.ID,
	)

	purchase.Status = status

	if err := db.Create(
		&purchase,
	).Error; err != nil {

		t.Fatal(err)
	}

	return &purchase
}
