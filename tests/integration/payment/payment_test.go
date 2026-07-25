package payment

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"
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

	err = helpers.CleanDatabase(
		db,
	)

	if err != nil {
		t.Fatal(err)
	}

	executor :=
		database.NewGormExecutor(
			db,
		)

	// -----------------------
	// Create fixtures
	// -----------------------

	category := fixtures.Category()

	if err := db.Create(
		&category,
	).Error; err != nil {

		t.Fatal(err)
	}

	user := fixtures.User()

	if err := db.Create(
		&user,
	).Error; err != nil {

		t.Fatal(err)
	}

	party := fixtures.PartyWithOrganizer(
		user.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(
		&party,
	).Error; err != nil {

		t.Fatal(err)
	}

	purchase :=
		helpers.CreatePurchase(
			t,
			db,
			&user,
			&party,
			enum.StatusPending,
		)

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

	paymentEventRepository :=
		repository.NewPaymentEventRepository(
			executor,
		)

	paymentService :=
		service.NewPaymentService(
			purchaseService,
			ticketService,
			fakeGateway,
			paymentEventRepository,
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

	if updated.Status != enum.StatusPending {

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
		enum.StatusPaid,
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
		enum.StatusPending,
	)

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
		enum.StatusPending,
	)

	_, err = paymentService.CreateCheckout(
		context.Background(),
		purchase.ID,
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

	if updatedPurchase.Status != enum.StatusPending {

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
		enum.StatusPending,
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

	if result.Status != enum.StatusPaid {

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

	if updatedPurchase.Status != enum.StatusPaid {

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
		enum.StatusPaid,
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

	if result.Status != enum.StatusPaid {

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
		enum.StatusPending,
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

	if updatedPurchase.Status != enum.StatusPaid {

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
		enum.StatusPending,
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
