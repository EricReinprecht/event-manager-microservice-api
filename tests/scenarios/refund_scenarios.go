package scenarios

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/service"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

type RefundScenario struct {
	Organizer models.User

	Actor models.User

	Category models.Category

	Party models.Party

	Purchase *models.Purchase

	Role enum.PartyRole
}

func CreateRefundScenario(
	t *testing.T,
	db *gorm.DB,
	role enum.PartyRole,
) RefundScenario {

	// -------------------------
	// Organizer
	// -------------------------

	organizer := fixtures.User()

	if err := db.Create(
		&organizer,
	).Error; err != nil {

		t.Fatal(err)
	}

	// -------------------------
	// Actor who performs refund
	// -------------------------

	actor := fixtures.User()

	if err := db.Create(
		&actor,
	).Error; err != nil {

		t.Fatal(err)
	}

	// -------------------------
	// Category
	// -------------------------

	category := fixtures.Category()

	if err := db.Create(
		&category,
	).Error; err != nil {

		t.Fatal(err)
	}

	// -------------------------
	// Party
	// -------------------------

	party := fixtures.PartyWithOrganizer(
		organizer.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(
		&party,
	).Error; err != nil {

		t.Fatal(err)
	}

	// -------------------------
	// Add role if needed
	// Organizer is already handled
	// by HasRole()
	// -------------------------

	if role != enum.RoleOrganizer {

		member := models.PartyMember{

			ID: uuid.New(),

			PartyID: party.ID,

			UserID: actor.ID,
		}

		if err := db.Create(
			&member,
		).Error; err != nil {

			t.Fatal(err)
		}

		memberRole := models.PartyMemberRole{

			ID: uuid.New(),

			PartyMemberID: member.ID,

			Role: role,
		}

		if err := db.Create(
			&memberRole,
		).Error; err != nil {

			t.Fatal(err)
		}
	}

	// -------------------------
	// Purchase
	// -------------------------

	purchase := helpers.CreatePurchase(
		t,
		db,
		&actor,
		&party,
		enum.PurchaseStatusPaid,
	)

	purchase.PaymentProvider = "paypal"

	purchase.PaymentID = "PAYPAL-REFUND-TEST"

	if err := db.Save(
		&purchase,
	).Error; err != nil {

		t.Fatal(err)
	}

	return RefundScenario{

		Organizer: organizer,

		Actor: actor,

		Category: category,

		Party: party,

		Purchase: purchase,

		Role: role,
	}
}

func RefundPurchase(
	t *testing.T,
	paymentService *service.PaymentService,
	purchaseID uuid.UUID,
	userID uuid.UUID,
) {

	err := paymentService.RefundPayment(
		context.Background(),
		purchaseID,
		userID,
	)

	if err != nil {

		t.Fatal(err)
	}
}
