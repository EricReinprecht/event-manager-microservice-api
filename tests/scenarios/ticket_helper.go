package scenarios

import (
	"testing"
	"time"

	"github.com/google/uuid"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"gorm.io/gorm"
)

func CreateAccessWindow(
	t *testing.T,
	db *gorm.DB,
	categoryID uuid.UUID,
	start time.Time,
	end time.Time,
) appModels.TicketAccessWindow {

	window := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: categoryID,

		StartsAt: start,

		EndsAt: end,
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	return window
}

func CreateNonVerificationCategory(
	t *testing.T,
	db *gorm.DB,
	partyID uuid.UUID,
) appModels.TicketCategory {

	category := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "Standard",

		Price: 50,

		Capacity: 100,

		PartyID: partyID,

		RequiresVerification: false,
	}

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	return category
}

func CreateTicket(
	t *testing.T,
	db *gorm.DB,
	categoryID uuid.UUID,
	userID uuid.UUID,
	partyID uuid.UUID,
) appModels.Ticket {

	purchase := fixtures.Purchase(
		userID,
		partyID,
	)

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = categoryID
	ticket.UserID = userID
	ticket.PurchaseID = purchase.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	return ticket
}

func CreateTestPurchase(
	t *testing.T,
	db *gorm.DB,
	userID uuid.UUID,
	partyID uuid.UUID,
) appModels.Purchase {

	t.Helper()

	purchase := appModels.Purchase{

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

func ChangeTicketOwner(
	t *testing.T,
	db *gorm.DB,
	ticket *appModels.Ticket,
	userID uuid.UUID,
) {

	if err := db.
		Model(ticket).
		Update(
			"user_id",
			userID,
		).
		Error; err != nil {

		t.Fatal(err)
	}
}
