package payment

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestPaymentCreateCheckoutSuccess(
	t *testing.T,
) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	// -----------------------
	// Create fixtures
	// -----------------------

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.PartyWithOrganizer(
		user.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	// -----------------------
	// Add purchase item
	// -----------------------

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	ticketCategory.Price = 1000

	if err := db.Create(
		&ticketCategory,
	).Error; err != nil {

		t.Fatal(err)
	}

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

	// -----------------------
	// Reload purchase
	// -----------------------

	if err := db.
		Preload("Items").
		First(
			&purchase,
			purchase.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// -----------------------
	// Services
	// -----------------------

	purchaseService :=
		helpers.NewPurchaseService(
			db,
		)

	ticketService :=
		helpers.NewTicketService(
			db,
		)

	fakeGateway :=
		&helpers.FakePaymentGateway{
			Order: &paypal.Order{
				ID:          "ORDER123",
				ApprovalURL: "https://paypal.test/approve",
			},
		}

	paymentService :=
		helpers.NewPaymentService(
			executor,
			purchaseService,
			ticketService,
			fakeGateway,
		)

	// -----------------------
	// Execute
	// -----------------------

	url, err :=
		paymentService.CreateCheckout(
			context.Background(),
			purchase.ID,
		)

	if err != nil {

		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	// -----------------------
	// Assertions
	// -----------------------

	if !fakeGateway.CreateOrderCalled {

		t.Fatal(
			"expected CreateOrder to be called",
		)
	}

	if url != "https://paypal.test/approve" {

		t.Fatalf(
			"unexpected approval url: %s",
			url,
		)
	}

	updated, err :=
		purchaseService.GetPurchase(
			context.Background(),
			purchase.ID,
		)

	if err != nil {
		t.Fatal(err)
	}

	if updated.PaymentProvider != "paypal" {

		t.Fatalf(
			"expected payment provider paypal, got %s",
			updated.PaymentProvider,
		)
	}

	if updated.PaymentID != "ORDER123" {

		t.Fatalf(
			"expected payment id ORDER123, got %s",
			updated.PaymentID,
		)
	}

	if updated.Status != enum.PurchaseStatusPending {

		t.Fatalf(
			"expected pending status, got %s",
			updated.Status,
		)
	}
}

func TestPaymentCreateCheckoutInvalidPurchase(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{
		Order: &paypal.Order{
			ID:          "ORDER-123",
			ApprovalURL: "https://paypal.test/approve",
		},
	}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	_, err = paymentService.CreateCheckout(
		context.Background(),
		uuid.New(), // non existing purchase ID
	)

	if err == nil {
		t.Fatal("expected error for invalid purchase")
	}

	if fakeGateway.CreateOrderCalled {

		t.Fatal(
			"paypal gateway should not be called for invalid purchase",
		)
	}
}

func TestPaymentCreateCheckoutAlreadyPaid(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{
		Order: &paypal.Order{
			ID:          "ORDER-123",
			ApprovalURL: "https://paypal.test/approve",
		},
	}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

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

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPaid,
	)

	_, err = paymentService.CreateCheckout(
		context.Background(),
		purchase.ID,
	)

	if err == nil {
		t.Fatal(
			"expected error when purchase is already paid",
		)
	}

	if err.Error() != "purchase already paid" {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if fakeGateway.CreateOrderCalled {

		t.Fatal(
			"paypal should not be called for already paid purchase",
		)
	}
}

func TestPaymentCreateCheckoutCreatesPayPalOrder(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{
		Order: &paypal.Order{
			ID:          "ORDER-123",
			ApprovalURL: "https://paypal.test/approve",
		},
	}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

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

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	// -----------------------
	// Add purchase item
	// -----------------------

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	ticketCategory.Price = 1000

	if err := db.Create(
		&ticketCategory,
	).Error; err != nil {
		t.Fatal(err)
	}

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

	// -----------------------
	// Execute
	// -----------------------

	url, err := paymentService.CreateCheckout(
		context.Background(),
		purchase.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if url != "https://paypal.test/approve" {

		t.Fatalf(
			"unexpected approval url: %s",
			url,
		)
	}

	if !fakeGateway.CreateOrderCalled {

		t.Fatal(
			"paypal create order was not called",
		)
	}

	var updatedPurchase models.Purchase

	err = db.
		First(
			&updatedPurchase,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.PaymentProvider != "paypal" {

		t.Fatalf(
			"expected paypal provider, got %s",
			updatedPurchase.PaymentProvider,
		)
	}

	if updatedPurchase.PaymentID != "ORDER-123" {

		t.Fatalf(
			"expected paypal order id ORDER-123, got %s",
			updatedPurchase.PaymentID,
		)
	}
}

func TestPaymentCreateCheckoutSavesPaymentID(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{
		Order: &paypal.Order{
			ID:          "PAYPAL-ORDER-456",
			ApprovalURL: "https://paypal.test/approve",
		},
	}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

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

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	// -----------------------
	// Add purchase item
	// -----------------------

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	ticketCategory.Price = 1000

	if err := db.Create(
		&ticketCategory,
	).Error; err != nil {
		t.Fatal(err)
	}

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

	// -----------------------
	// Create checkout
	// -----------------------

	_, err = paymentService.CreateCheckout(
		context.Background(),
		purchase.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	// -----------------------
	// Verify
	// -----------------------

	var updatedPurchase models.Purchase

	err = db.
		First(
			&updatedPurchase,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.PaymentID != "PAYPAL-ORDER-456" {

		t.Fatalf(
			"expected payment id PAYPAL-ORDER-456, got %s",
			updatedPurchase.PaymentID,
		)
	}

	if updatedPurchase.PaymentProvider != "paypal" {

		t.Fatalf(
			"expected payment provider paypal, got %s",
			updatedPurchase.PaymentProvider,
		)
	}

	if updatedPurchase.Status != enum.PurchaseStatusPending {

		t.Fatalf(
			"expected purchase to remain pending, got %s",
			updatedPurchase.Status,
		)
	}
}

func TestPaymentConfirmPaymentSuccess(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

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

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = "PAYPAL-ORDER-123"

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{
		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 1,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	purchase.Items = []models.PurchaseItem{
		item,
	}

	result, err := paymentService.ConfirmPayment(
		context.Background(),
		"PAYPAL-ORDER-123",
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected purchase status PAID, got %s",
			result.Status,
		)
	}

	var updatedPurchase models.Purchase

	err = db.
		First(
			&updatedPurchase,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected database purchase status PAID, got %s",
			updatedPurchase.Status,
		)
	}

	var tickets []models.Ticket

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(tickets) != 1 {

		t.Fatalf(
			"expected 1 generated ticket, got %d",
			len(tickets),
		)
	}
}

func TestPaymentConfirmPaymentPurchaseNotFound(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	_, err = paymentService.ConfirmPayment(
		context.Background(),
		"NON_EXISTING_PAYPAL_ORDER",
	)

	if err == nil {

		t.Fatal(
			"expected error when purchase was not found",
		)
	}
}

func TestPaymentConfirmPaymentAlreadyPaid(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

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

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPaid,
	)

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = "PAYPAL-ORDER-PAID"

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{
		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 1,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	result, err := paymentService.ConfirmPayment(
		context.Background(),
		"PAYPAL-ORDER-PAID",
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected PAID status, got %s",
			result.Status,
		)
	}

	var tickets []models.Ticket

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(tickets) != 0 {

		t.Fatalf(
			"expected no tickets to be generated again, got %d",
			len(tickets),
		)
	}
}

func TestPaymentConfirmPaymentGeneratesTickets(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

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

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = "PAYPAL-ORDER-GENERATE-TICKETS"

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{
		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 2,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	_, err = paymentService.ConfirmPayment(
		context.Background(),
		"PAYPAL-ORDER-GENERATE-TICKETS",
	)

	if err != nil {
		t.Fatal(err)
	}

	var updatedPurchase models.Purchase

	err = db.
		First(
			&updatedPurchase,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected purchase status PAID, got %s",
			updatedPurchase.Status,
		)
	}

	var tickets []models.Ticket

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(tickets) != 2 {

		t.Fatalf(
			"expected 2 generated tickets, got %d",
			len(tickets),
		)
	}

	for _, ticket := range tickets {

		if ticket.TicketCategoryID != ticketCategory.ID {

			t.Fatalf(
				"expected ticket category %s, got %s",
				ticketCategory.ID,
				ticket.TicketCategoryID,
			)
		}

		if ticket.Code == "" {

			t.Fatal(
				"expected generated ticket code",
			)
		}
	}
}

func TestPaymentConfirmPaymentDoesNotDuplicateTickets(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

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

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = "PAYPAL-ORDER-IDEMPOTENCY"

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{
		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 2,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	_, err = paymentService.ConfirmPayment(
		context.Background(),
		"PAYPAL-ORDER-IDEMPOTENCY",
	)

	if err != nil {
		t.Fatal(err)
	}

	_, err = paymentService.ConfirmPayment(
		context.Background(),
		"PAYPAL-ORDER-IDEMPOTENCY",
	)

	if err != nil {
		t.Fatal(err)
	}

	var tickets []models.Ticket

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(tickets) != 2 {

		t.Fatalf(
			"expected exactly 2 tickets after duplicate confirmation, got %d",
			len(tickets),
		)
	}
}

func TestPaymentConfirmPaymentInvalidPayPalOrder(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	_, err = paymentService.ConfirmPayment(
		context.Background(),
		"INVALID_PAYPAL_ORDER_ID",
	)

	if err == nil {

		t.Fatal(
			"expected invalid PayPal order error",
		)
	}
}

func TestPurchase_StatusChangesPendingToPaid(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

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

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	paymentID := "PAYPAL-STATUS-" + uuid.New().String()

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = paymentID

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// verify initial status

	if purchase.Status != enum.PurchaseStatusPending {

		t.Fatalf(
			"expected initial status PENDING, got %s",
			purchase.Status,
		)
	}

	_, err = paymentService.ConfirmPayment(
		context.Background(),
		paymentID,
	)

	if err != nil {

		t.Fatal(
			"confirm payment failed:",
			err,
		)
	}

	var updatedPurchase models.Purchase

	err = db.
		First(
			&updatedPurchase,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected status PAID, got %s",
			updatedPurchase.Status,
		)
	}
}

func TestPurchase_StatusDoesNotMoveBackwards(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

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

	paymentID := "PAYPAL-NO-BACKWARD-" + uuid.New().String()

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPaid,
	)

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = paymentID

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// Confirming an already paid purchase should not downgrade it

	result, err := paymentService.ConfirmPayment(
		context.Background(),
		paymentID,
	)

	if err != nil {
		t.Fatal(
			"confirm payment failed:",
			err,
		)
	}

	if result.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected returned status PAID, got %s",
			result.Status,
		)
	}

	var updatedPurchase models.Purchase

	err = db.
		First(
			&updatedPurchase,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected database status to remain PAID, got %s",
			updatedPurchase.Status,
		)
	}
}

func TestPayment_GeneratesCorrectNumberOfTickets(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

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

	paymentID := "PAYPAL-TICKETS-" + uuid.New().String()

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = paymentID

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	quantity := 3

	item := models.PurchaseItem{

		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: quantity,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	purchase.Items = []models.PurchaseItem{
		item,
	}

	_, err = paymentService.ConfirmPayment(
		context.Background(),
		paymentID,
	)

	if err != nil {
		t.Fatal(
			"confirm payment failed:",
			err,
		)
	}

	var tickets []models.Ticket

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(tickets) != quantity {

		t.Fatalf(
			"expected %d tickets, got %d",
			quantity,
			len(tickets),
		)
	}
}

func TestPayment_DoesNotGenerateDuplicateTickets(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

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

	paymentID := "PAYPAL-NO-DUPLICATE-TICKETS-" + uuid.New().String()

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = paymentID

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{

		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 2,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	purchase.Items = []models.PurchaseItem{
		item,
	}

	// First confirmation

	_, err = paymentService.ConfirmPayment(
		context.Background(),
		paymentID,
	)

	if err != nil {
		t.Fatal(err)
	}

	var tickets []models.Ticket

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	firstCount := len(tickets)

	if firstCount != 2 {

		t.Fatalf(
			"expected 2 tickets after first confirmation, got %d",
			firstCount,
		)
	}

	// Second confirmation

	_, err = paymentService.ConfirmPayment(
		context.Background(),
		paymentID,
	)

	if err != nil {
		t.Fatal(err)
	}

	tickets = nil

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	secondCount := len(tickets)

	if secondCount != firstCount {

		t.Fatalf(
			"expected ticket count to remain %d, got %d",
			firstCount,
			secondCount,
		)
	}
}

func TestPaymentConfirmPaymentConcurrentCalls(t *testing.T) {

	db, err := helpers.TestDatabase()
	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)
	if err != nil {
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

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{

		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 1,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	// Fake PayPal order id

	paymentID := "PAYPAL-CONCURRENT-TEST"

	purchase.PaymentID = paymentID

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// ----------------------------
	// Concurrent confirmation
	// ----------------------------

	var wg sync.WaitGroup

	errors := make(chan error, 2)

	wg.Add(2)

	for i := 0; i < 2; i++ {

		go func() {

			defer wg.Done()

			_, err := paymentService.ConfirmPayment(
				context.Background(),
				paymentID,
			)

			errors <- err

		}()

	}

	wg.Wait()

	close(errors)

	// ----------------------------
	// Validate result
	// ----------------------------

	for err := range errors {

		if err != nil {

			t.Log(
				"concurrent confirm returned:",
				err,
			)
		}
	}

	var tickets []models.Ticket

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(tickets) != 1 {

		t.Fatalf(
			"expected exactly 1 ticket, got %d",
			len(tickets),
		)
	}

	var updatedPurchase models.Purchase

	err = db.
		First(
			&updatedPurchase,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected PAID got %s",
			updatedPurchase.Status,
		)
	}
}

func TestConcurrentConfirmPayment(t *testing.T) {

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

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	ticketCategory := fixtures.TicketCategory(party.ID)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{
		ID:               uuid.New(),
		PurchaseID:       purchase.ID,
		TicketCategoryID: ticketCategory.ID,
		Quantity:         1,
		UnitPrice:        ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = "PAYPAL-CONCURRENT"

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// ----------------------------
	// Execute concurrently
	// ----------------------------

	const workers = 10

	var wg sync.WaitGroup

	wg.Add(workers)

	for i := 0; i < workers; i++ {

		go func() {

			defer wg.Done()

			_, _ = paymentService.ConfirmPayment(
				context.Background(),
				purchase.PaymentID,
			)

		}()

	}

	wg.Wait()

	// ----------------------------
	// Verify purchase
	// ----------------------------

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		purchase.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected purchase to be PAID, got %s",
			updated.Status,
		)
	}

	// ----------------------------
	// Verify tickets
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
			"expected exactly 1 ticket after concurrent confirmation, got %d",
			len(tickets),
		)
	}
}

func TestPaymentCaptureFailureRollback(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

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

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	ticketCategory := fixtures.TicketCategory(party.ID)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{
		ID:               uuid.New(),
		PurchaseID:       purchase.ID,
		TicketCategoryID: ticketCategory.ID,
		Quantity:         2,
		UnitPrice:        ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = "INVALID-ORDER"

	if err := db.Save(&purchase).Error; err != nil {
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
	// Execute
	// ----------------------------

	_, err = paymentService.CapturePayment(
		context.Background(),
		purchase.PaymentID,
	)
	if err == nil {
		t.Fatal("expected capture to fail")
	}

	// ----------------------------
	// Purchase should still be pending
	// ----------------------------

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		purchase.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.PurchaseStatusPending {

		t.Fatalf(
			"expected purchase to remain PENDING, got %s",
			updated.Status,
		)
	}

	// ----------------------------
	// No tickets generated
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

	if len(tickets) != 0 {

		t.Fatalf(
			"expected 0 tickets after failed capture, got %d",
			len(tickets),
		)
	}
}

func TestPaymentSupportsMultipleTicketCategories(t *testing.T) {

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

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	vip := fixtures.TicketCategory(party.ID)
	vip.Name = "VIP"

	if err := db.Create(&vip).Error; err != nil {
		t.Fatal(err)
	}

	standard := fixtures.TicketCategory(party.ID)
	standard.Name = "Standard"

	if err := db.Create(&standard).Error; err != nil {
		t.Fatal(err)
	}

	items := []models.PurchaseItem{
		{
			ID:               uuid.New(),
			PurchaseID:       purchase.ID,
			TicketCategoryID: vip.ID,
			Quantity:         2,
			UnitPrice:        vip.Price,
		},
		{
			ID:               uuid.New(),
			PurchaseID:       purchase.ID,
			TicketCategoryID: standard.ID,
			Quantity:         3,
			UnitPrice:        standard.Price,
		},
	}

	if err := db.Create(&items).Error; err != nil {
		t.Fatal(err)
	}

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = "MULTI-CATEGORY"

	if err := db.Save(&purchase).Error; err != nil {
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
	// Verify tickets
	// ----------------------------

	var tickets []models.Ticket

	if err := db.
		Where("user_id = ?", user.ID).
		Find(&tickets).
		Error; err != nil {

		t.Fatal(err)
	}

	if len(tickets) != 5 {

		t.Fatalf(
			"expected 5 tickets, got %d",
			len(tickets),
		)
	}

	var vipCount int
	var standardCount int

	for _, ticket := range tickets {

		switch ticket.TicketCategoryID {

		case vip.ID:
			vipCount++

		case standard.ID:
			standardCount++
		}
	}

	if vipCount != 2 {

		t.Fatalf(
			"expected 2 VIP tickets, got %d",
			vipCount,
		)
	}

	if standardCount != 3 {

		t.Fatalf(
			"expected 3 Standard tickets, got %d",
			standardCount,
		)
	}
}
