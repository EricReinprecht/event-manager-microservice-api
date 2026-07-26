package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/requests"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func setupPurchaseTest(t *testing.T) (
	*appModels.User,
	*appModels.Party,
) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.PartyWithOrganizer(
		user.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	return &user, &party
}

func TestUserCanCreatePendingPurchase(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	user, party := setupPurchaseTest(t)

	service := helpers.NewPurchaseService(db)

	category := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	purchase, err := service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: category.ID,
				Quantity:         1,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if purchase.Status != enum.PurchaseStatusPending {

		t.Fatalf(
			"expected pending purchase got %s",
			purchase.Status,
		)
	}

	if len(purchase.Items) != 1 {

		t.Fatalf(
			"expected one purchase item, got %d",
			len(purchase.Items),
		)
	}

	if purchase.Items[0].UnitPrice != category.Price {

		t.Fatalf(
			"expected snapshot price %d got %d",
			category.Price,
			purchase.Items[0].UnitPrice,
		)
	}
}

func TestPurchaseStoresPriceSnapshot(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	user, party := setupPurchaseTest(t)

	service := helpers.NewPurchaseService(db)

	category := fixtures.TicketCategory(
		party.ID,
	)

	category.Price = 5000

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	purchase, err := service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: category.ID,
				Quantity:         2,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if purchase.Items[0].UnitPrice != 5000 {

		t.Fatalf(
			"expected snapshot price 5000 got %d",
			purchase.Items[0].UnitPrice,
		)
	}
}

func TestPurchaseRejectsInvalidCategory(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	user, party := setupPurchaseTest(t)

	service := helpers.NewPurchaseService(db)

	_, err = service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: uuid.New(),
				Quantity:         1,
			},
		},
	)

	if err == nil {

		t.Fatal(
			"expected invalid category error",
		)
	}
}

func TestPurchaseCalculatesTotal(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	user, party := setupPurchaseTest(t)

	service := helpers.NewPurchaseService(db)

	category := fixtures.TicketCategory(
		party.ID,
	)

	category.Price = 2500

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	purchase, err := service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: category.ID,
				Quantity:         4,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	expected := int64(10000)

	if purchase.TotalPrice != expected {

		t.Fatalf(
			"expected total %d got %d",
			expected,
			purchase.TotalPrice,
		)
	}
}

func TestCannotPurchaseTicketCategoryFromAnotherParty(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	user, party := setupPurchaseTest(t)

	otherCategory := fixtures.Category()
	otherCategory.Name = "Concert"

	if err := db.Create(&otherCategory).Error; err != nil {
		t.Fatal(err)
	}

	otherParty := fixtures.PartyWithOrganizer(
		user.ID,
	)

	otherParty.CategoryID = otherCategory.ID

	if err := db.Create(&otherParty).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.TicketCategory(
		otherParty.ID,
	)

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	service := helpers.NewPurchaseService(db)

	_, err = service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: category.ID,
				Quantity:         1,
			},
		},
	)

	if err == nil {

		t.Fatal(
			"expected category from another party to fail",
		)
	}

	if err != appErrors.ErrTicketCategoryNotFound {

		t.Fatalf(
			"expected ErrTicketCategoryNotFound got %v",
			err,
		)
	}
}

func TestPurchaseCanBeMarkedPaid(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	user, party := setupPurchaseTest(t)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	paymentService := helpers.NewPaymentService(
		database.NewGormExecutor(db),
		purchaseService,
		ticketService,
		helpers.NewPayPalClient(),
	)

	category := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	purchase, err := purchaseService.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: category.ID,
				Quantity:         1,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	paymentID := "PAYPAL-ORDER-123"

	err = purchaseService.AttachPayment(
		context.Background(),
		purchase.ID,
		"paypal",
		paymentID,
	)

	if err != nil {
		t.Fatal(err)
	}

	updatedPurchase, err := paymentService.ConfirmPayment(
		context.Background(),
		paymentID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected purchase status PAID got %s",
			updatedPurchase.Status,
		)
	}

	if updatedPurchase.PaymentProvider != "paypal" {

		t.Fatalf(
			"expected payment provider paypal got %s",
			updatedPurchase.PaymentProvider,
		)
	}

	if updatedPurchase.PaymentID != paymentID {

		t.Fatalf(
			"expected payment id %s got %s",
			paymentID,
			updatedPurchase.PaymentID,
		)
	}
}

func TestPaymentUsesPurchasePriceSnapshot(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	paymentService := helpers.NewPaymentService(
		database.NewGormExecutor(db),
		purchaseService,
		ticketService,
		helpers.NewPayPalClient(),
	)

	// ----------------------------
	// Setup
	// ----------------------------

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.Party()

	party.CategoryID = category.ID
	party.OrganizerID = user.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	ticketCategory.Price = 5000

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// ----------------------------
	// Create purchase with snapshot
	// ----------------------------

	purchase, err := purchaseService.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: ticketCategory.ID,
				Quantity:         1,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	purchase.ExpiresAt = time.Now().Add(
		30 * time.Minute,
	)

	purchase.PaymentProvider = "paypal"

	purchase.PaymentID = "PRICE-SNAPSHOT"

	if err := db.Save(
		purchase,
	).Error; err != nil {

		t.Fatal(err)
	}

	// ----------------------------
	// Change ticket price AFTER purchase
	// ----------------------------

	ticketCategory.Price = 9000

	if err := db.Save(
		&ticketCategory,
	).Error; err != nil {

		t.Fatal(err)
	}

	// ----------------------------
	// Confirm payment
	// ----------------------------

	_, err = paymentService.ConfirmPayment(
		context.Background(),
		purchase.PaymentID,
	)

	if err != nil {
		t.Fatal(err)
	}

	// ----------------------------
	// Reload purchase
	// ----------------------------

	var updated models.Purchase

	if err := db.
		Preload("Items").
		First(
			&updated,
			"id = ?",
			purchase.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	if len(updated.Items) != 1 {

		t.Fatalf(
			"expected 1 purchase item, got %d",
			len(updated.Items),
		)
	}

	if updated.Items[0].UnitPrice != 5000 {

		t.Fatalf(
			"expected purchase snapshot price 5000, got %d",
			updated.Items[0].UnitPrice,
		)
	}

	// ----------------------------
	// Verify ticket created
	// ----------------------------

	var tickets []models.Ticket

	if err := db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error; err != nil {

		t.Fatal(err)
	}

	if len(tickets) != 1 {

		t.Fatalf(
			"expected 1 ticket, got %d",
			len(tickets),
		)
	}
}

func TestPurchaseRejectsSoldOutTicketCategory(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	user, party := setupPurchaseTest(t)

	service := helpers.NewPurchaseService(db)

	category := fixtures.TicketCategory(
		party.ID,
	)

	category.Capacity = 2

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// Create a purchase because tickets require a valid PurchaseID

	purchase := scenarios.CreatePurchase(
		t,
		db,
		user.ID,
		party.ID,
	)

	// Simulate that all tickets were already sold

	tickets := []appModels.Ticket{
		{
			ID: uuid.New(),

			Code: "SOLD001",

			Status: enum.TicketStatusActive,

			UserID: user.ID,

			TicketCategoryID: category.ID,

			PurchaseID: purchase.ID,
		},
		{
			ID: uuid.New(),

			Code: "SOLD002",

			Status: enum.TicketStatusActive,

			UserID: user.ID,

			TicketCategoryID: category.ID,

			PurchaseID: purchase.ID,
		},
	}

	if err := db.Create(&tickets).Error; err != nil {
		t.Fatal(err)
	}

	_, err = service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: category.ID,
				Quantity:         1,
			},
		},
	)

	if err == nil {
		t.Fatal(
			"expected sold-out category to be rejected",
		)
	}

	if !errors.Is(
		err,
		appErrors.ErrTicketSoldOut,
	) {
		t.Fatalf(
			"expected ErrTicketSoldOut, got %v",
			err,
		)
	}
}

func TestPurchaseTransactionRollbackOnFailure(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	user, party := setupPurchaseTest(t)

	service := helpers.NewPurchaseService(db)

	category := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	_, err = service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{

			{
				TicketCategoryID: category.ID,
				Quantity:         1,
			},

			{
				// invalid category forces failure
				TicketCategoryID: uuid.New(),
				Quantity:         1,
			},
		},
	)

	if err == nil {

		t.Fatal(
			"expected purchase creation to fail",
		)
	}

	// ----------------------------
	// Verify rollback
	// ----------------------------

	var purchases []appModels.Purchase

	if err := db.Find(
		&purchases,
	).Error; err != nil {

		t.Fatal(err)
	}

	if len(purchases) != 0 {

		t.Fatalf(
			"expected zero purchases after rollback, got %d",
			len(purchases),
		)
	}

	var items []appModels.PurchaseItem

	if err := db.Find(
		&items,
	).Error; err != nil {

		t.Fatal(err)
	}

	if len(items) != 0 {

		t.Fatalf(
			"expected zero purchase items after rollback, got %d",
			len(items),
		)
	}
}
