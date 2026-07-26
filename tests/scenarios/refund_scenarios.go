package scenarios

import (
	"context"
	"testing"
	"time"

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

	TicketCategory models.TicketCategory

	RefundPolicy models.RefundPolicy

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

	if role == enum.RoleOrganizer {
		actor = organizer
	}

	if role != enum.RoleOrganizer {

		if err := db.Create(
			&actor,
		).Error; err != nil {

			t.Fatal(err)
		}
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
	// -------------------------

	if role != enum.RoleOrganizer {

		member := models.PartyMember{
			ID:      uuid.New(),
			PartyID: party.ID,
			UserID:  actor.ID,
		}

		if err := db.Create(
			&member,
		).Error; err != nil {

			t.Fatal(err)
		}

		memberRole := models.PartyMemberRole{
			ID:            uuid.New(),
			PartyMemberID: member.ID,
			Role:          role,
		}

		if err := db.Create(
			&memberRole,
		).Error; err != nil {

			t.Fatal(err)
		}
	}

	// -------------------------
	// Ticket Category
	// -------------------------

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	ticketCategory.Price = 1000

	if err := db.Create(
		&ticketCategory,
	).Error; err != nil {

		t.Fatal(err)
	}

	// -------------------------
	// Refund Policy
	// -------------------------

	refundPolicy := models.RefundPolicy{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		Until: time.Now().AddDate(
			0,
			0,
			30,
		),

		Percentage: 100,
	}

	if err := db.Create(
		&refundPolicy,
	).Error; err != nil {

		t.Fatal(err)
	}

	// reload ticket category with refund policy
	if err := db.
		Preload("RefundPolicy").
		First(
			&ticketCategory,
			ticketCategory.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// -------------------------
	// Attach Refund Policy
	// -------------------------

	ticketCategory.RefundPolicyID = &refundPolicy.ID

	if err := db.Save(
		&ticketCategory,
	).Error; err != nil {

		t.Fatal(err)
	}

	// reload ticket category with policy
	if err := db.
		Preload("RefundPolicy").
		First(
			&ticketCategory,
			ticketCategory.ID,
		).
		Error; err != nil {

		t.Fatal(err)
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

	// -------------------------
	// Purchase Item
	// -------------------------

	item := models.PurchaseItem{

		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 1,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(
		&item,
	).Error; err != nil {

		t.Fatal(err)
	}

	// reload purchase with relations
	if err := db.
		Preload("Items.TicketCategory.RefundPolicy").
		First(
			&purchase,
			purchase.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	return RefundScenario{

		Organizer: organizer,

		Actor: actor,

		Category: category,

		Party: party,

		Purchase: purchase,

		TicketCategory: ticketCategory,

		RefundPolicy: refundPolicy,

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
