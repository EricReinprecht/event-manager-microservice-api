package helpers

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"

	appModels "github.com/reinp/event-platform/backend/internal/models"
)

func NewPurchaseService(db *gorm.DB) *service.PurchaseService {

	executor := database.NewGormExecutor(db)

	purchaseRepository := repository.NewPurchaseRepository(
		executor,
	)

	ticketRepository := repository.NewTicketRepository(
		executor,
	)

	return service.NewPurchaseService(
		purchaseRepository,
		ticketRepository,
	)
}

func CreatePurchase(
	t *testing.T,
	db *gorm.DB,
	user *models.User,
	party *models.Party,
	status enum.PurchaseStatus,
) *models.Purchase {

	purchase := models.Purchase{

		ID: uuid.New(),

		UserID: user.ID,

		PartyID: party.ID,

		Status: status,

		ExpiresAt: time.Now().Add(
			30 * time.Minute,
		),
	}

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	return &purchase
}

func CreateTestPurchase(
	db *gorm.DB,
	userID uuid.UUID,
	partyID uuid.UUID,
) appModels.Purchase {

	purchase := appModels.Purchase{
		ID: uuid.New(),

		UserID: userID,

		PartyID: partyID,

		Status: enum.PurchaseStatusPaid,
	}

	if err := db.Create(&purchase).Error; err != nil {
		panic(err)
	}

	return purchase
}
